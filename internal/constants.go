package internal

type Register struct {
	name     string `form:"name"`
	email    string `form:"email"`
	phone    string `form:"phone"`
	password string `form:"password"`
}

type RegisterCompany struct {
	registerPointer *Register
	companyName     string `form:"company_name"`
	post            string `form:"post"`
	taxpayerNumber  string `form:"taxpayer_number"`
	address         string `form:"address"`
	serviceType     string `form:"service_type"`
}
