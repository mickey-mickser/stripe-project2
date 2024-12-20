package api

type User struct {
	Id       int     `json:"-" db:"id" gorm:"primaryKey;autoIncrement"` // ID, скрыт для клиента
	Name     string  `json:"name" binding:"required"`                   // Имя пользователя
	Username string  `json:"username" binding:"required"`               // Уникальное имя
	Balance  float64 `json:"balance" gorm:"default:0"`                  // Баланс
	Password string  `json:"password" binding:"required"`               // Пароль
}
type PaymentSession struct {
	ID        int64   `gorm:"primaryKey;autoIncrement" json:"id"`                         // Автоинкрементный идентификатор
	SessionID string  `gorm:"type:varchar(255);unique;not null" json:"session_id"`        // Уникальный session_id
	Username  string  `gorm:"type:varchar(255);not null" json:"username"`                 // Имя пользователя
	Amount    float64 `gorm:"type:numeric(10,2);not null" json:"amount"`                  // Сумма (с точностью до 2 знаков)
	Status    string  `gorm:"type:varchar(50);default:'pending'" json:"status"`           // Статус сессии
	CreatedAt string  `gorm:"type:timestamp;default:current_timestamp" json:"created_at"` // Время создания записи
}

//
//type User struct {
//	Name string `gorm:"<-:create"` // разрешены чтение и создание
//	Name string `gorm:"<-:update"` // разрешены чтение и обновление
//	Name string `gorm:"<-"`        // разрешены чтение и запись (создание и обновление)
//	Name string `gorm:"<-:false"`  // разрешено чтение, отключена запись
//	Name string `gorm:"->"`        // только чтение (отключена запись если она не настроена)
//	Name string `gorm:"->;<-:create"` // разрешены чтение и создание
//	Name string `gorm:"->:false;<-:create"` // только создание (отключено чтение из базы данных)
//	Name string `gorm:"-"`            // игнорируйте это поле при записи и чтении с помощью struct
//	Name string `gorm:"-:all"`        // игнорируйте это поле при записи, чтении и миграции с помощью struct
//	Name string `gorm:"-:migration"`  // игнорируйте это поле при миграции с помощью struct
//}
