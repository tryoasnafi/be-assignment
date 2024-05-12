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

func Init(userService *user.UserService) error {
	apiBasePath := "api/auth"
	websiteBasePath := "/auth"
	err := supertokens.Init(supertokens.TypeInput{
		Debug: true,
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
