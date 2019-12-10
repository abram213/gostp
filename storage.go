package gostp

type regexAndDescription struct {
	regex       string
	description string
}

var functionsMap map[string]interface{}

var regexMap map[string]regexAndDescription

// InitRegex - initialize all regexes which needed to precheck values
func InitRegex() {
	regexMap = make(map[string]regexAndDescription)
	// Checks if email address
	regexMap["email"] = regexAndDescription{"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", "email isn't valid"}
	// Checks if valid russian phone
	regexMap["russianPhone"] = regexAndDescription{`^(\+7|7|8)?[\s\-]?\(?[489][0-9]{2}\)?[\s\-]?[0-9]{3}[\s\-]?[0-9]{2}[\s\-]?[0-9]{2}$`, "phone isn't valid"}
	// Checks if password is more than 6 symbols
	regexMap["password"] = regexAndDescription{`^.{6,}$`, "password is less than 6 symbols"}

	// Functions Map initialization
	functionsMap = map[string]interface{}{
		"hashpwd": HashPassword,
	}
}
