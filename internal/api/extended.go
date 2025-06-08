package api

// Расширенные структуры для карточек услуг
type TokenUpdateCard struct {
	*TokenAccess
	CardID uint `json:"card_id"`
	Card   struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Location    string  `json:"location"`
		Price       float64 `json:"price"`
	} `json:"card"`
}

type ExtendedCardResponse struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Location    string  `json:"location"`
	Price       float64 `json:"price"`
	IsActive    bool    `json:"is_active"`
	CompanyID   uint    `json:"company_id"`
	Company     struct {
		ID          uint    `json:"id"`
		CompanyName string  `json:"company_name"`
		Stars       float64 `json:"stars"`
		ReviewCount int     `json:"review_count"`
		Photo       string  `json:"photo"`
	} `json:"company"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Структуры для заказов
type TokenCreateOrder struct {
	*TokenAccess
	Order struct {
		CompanyID   uint   `json:"company_id"`
		CardID      uint   `json:"card_id"`
		Description string `json:"description"`
	} `json:"order"`
}

type TokenOrderAction struct {
	*TokenAccess
	OrderID uint `json:"order_id"`
}

type OrderResponse struct {
	ID                uint    `json:"id"`
	ClientID          uint    `json:"client_id"`
	CompanyID         uint    `json:"company_id"`
	CardID            uint    `json:"card_id"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
	PaymentStatus     string  `json:"payment_status"`
	Description       string  `json:"description"`
	WorkerCompleteURL string  `json:"worker_complete_url,omitempty"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	CompletedAt       *string `json:"completed_at,omitempty"`
	Client            struct {
		ID       uint   `json:"id"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Photo    string `json:"photo"`
	} `json:"client,omitempty"`
	Company struct {
		ID          uint    `json:"id"`
		CompanyName string  `json:"company_name"`
		Stars       float64 `json:"stars"`
		Photo       string  `json:"photo"`
	} `json:"company,omitempty"`
	Card struct {
		ID          uint   `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"card,omitempty"`
}

// Структуры для баланса
type TokenDepositBalance struct {
	*TokenAccess
	Amount float64 `json:"amount"`
}

type TokenWithdrawBalance struct {
	*TokenAccess
	Amount float64 `json:"amount"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

type BalanceTransactionResponse struct {
	ID          uint    `json:"id"`
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"`
	Status      string  `json:"status"`
	Description string  `json:"description"`
	CreatedAt   string  `json:"created_at"`
}

// Структуры для отзывов
type TokenCreateReview struct {
	*TokenAccess
	Review struct {
		OrderID   uint   `json:"order_id"`
		CompanyID uint   `json:"company_id"`
		Rating    int    `json:"rating"`
		Comment   string `json:"comment"`
	} `json:"review"`
}

type ReviewResponse struct {
	ID        uint   `json:"id"`
	ClientID  uint   `json:"client_id"`
	CompanyID uint   `json:"company_id"`
	OrderID   uint   `json:"order_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
	Client    struct {
		ID       uint   `json:"id"`
		FullName string `json:"full_name"`
		Photo    string `json:"photo"`
	} `json:"client"`
}

// Структуры для работников
type WorkerCompleteRequest struct {
	Token string `json:"token"`
}

// Общие структуры для ответов
type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int         `json:"total"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Расширенные структуры для токен доступа с дополнительными полями
type ExtendedTokenAccess struct {
	*TokenAccess
	User struct {
		Login struct {
			Token string `json:"token"`
		} `json:"login"`
		ClientID  uint `json:"client_id,omitempty"`
		CompanyID uint `json:"company_id,omitempty"`
		UserType  string `json:"user_type"` // "client" или "company"
	} `json:"user"`
}

// Обновляем TokenDeleteCard для корректного использования
type TokenDeleteCardFixed struct {
	*TokenAccess
	CardID uint `json:"card_id"`
}

// Обновляем TokenCreateCard для расширенного функционала
type TokenCreateCardExtended struct {
	*TokenAccess
	Card struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Category    string  `json:"category"`
		Location    string  `json:"location"`
		Price       float64 `json:"price"`
	} `json:"card"`
}
