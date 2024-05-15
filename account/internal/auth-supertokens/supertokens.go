package auth_supertokens

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/jwt"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/tryoasnafi/be-assignment/account/internal/user"
	"gorm.io/datatypes"
)

type SupertokensAPIFunc = func(originalImplementation epmodels.APIInterface) epmodels.APIInterface

// Auth Supertokens
//	@Summary	Sign in to get token
//	@Schemes
//	@Description	get token from header fields (check docs for more details request example)
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	AuthResponse
//	@Router			/auth/signin [post]
//	@Param			user	body	AuthRequest	true	"auth signin"
// {
//   "formFields": [
//     {
//       "id": "email",
//       "value": "asnafi@gmail.com"
//     },
//     {
//       "id": "password",
//       "value": "helloworld1"
//     }
//   ]
// }
func Init(userService *user.UserService) error {
	apiBasePath := "api/auth"
	websiteBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: os.Getenv("SUPERTOKENS_URI"),
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "Simple Bank API",
			APIDomain:       "http://localhost:9090",
			WebsiteDomain:   "http://localhost:9090",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(&epmodels.TypeInput{
				Override: &epmodels.OverrideStruct{
					APIs: CustomSignupAPI(userService),
				},
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{ID: "first_name"},
						{ID: "last_name"},
						{ID: "dob"},
						{ID: "address"},
						{ID: "phone_number"},
					},
				},
			}),
			session.Init(nil),
			jwt.Init(nil),
		},
	})
	return err
}

// Auth Supertokens
//	@Summary	Sign up user
//	@Schemes
//	@Description	register user (check docs for more details request example)
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	AuthResponse
//	@Router			/auth/signup [post]
//	@Param			user	body	AuthRequest	true	"auth signup"
// {
//   "formFields": [
//     {
//       "id": "email",
//       "value": "asnafi1@gmail.com"
//     },
//     {
//       "id": "password",
//       "value": "helloworld1"
//     },
//      {
//       "id": "first_name",
//       "value": "asnafi"
//     },
//     {
//       "id": "last_name",
//       "value": "123"
//     },
//     {
//       "id": "dob",
//       "value": "2000-01-20"
//     },
//     {
//       "id": "phone_number",
//       "value": "+62817238123980"
//     },
//     {
//       "id": "address",
//       "value": "JL.OK"
//     }
//   ]
// }
func CustomSignupAPI(userService *user.UserService) SupertokensAPIFunc {
	return func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
		originalSignUpPOST := *originalImplementation.SignUpPOST

		(*originalImplementation.SignUpPOST) = func(
			formFields []epmodels.TypeFormField,
			tenantId string,
			options epmodels.APIOptions,
			userContext supertokens.UserContext,
		) (epmodels.SignUpPOSTResponse, error) {

			resp, err := originalSignUpPOST(formFields, tenantId, options, userContext)
			if err != nil {
				return epmodels.SignUpPOSTResponse{}, err
			}

			if resp.OK != nil {
				userId := resp.OK.User.ID
				// sign up successful
				u := &user.User{
					UUID: uuid.MustParse(userId),
				}
				for _, field := range formFields {
					switch field.ID {
					case "first_name":
						u.FirstName = field.Value
					case "last_name":
						u.LastName = field.Value
					case "email":
						u.Email = field.Value
					case "address":
						u.Address = field.Value
					case "dob":
						t, _ := time.Parse("2006-01-02", field.Value)
						u.DOB = datatypes.Date(t)
					case "phone_number":
						u.PhoneNumber = field.Value
					}
				}

				log.Println("Registered user", *u)
				userService.Create(u)
			}

			return resp, nil
		}

		return originalImplementation
	}
}
