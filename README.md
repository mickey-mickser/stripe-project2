# stripe_project2

На платформе есть пользователи  
У пользователей есть балансы в долларах
Нужно написать бэкенд, с помощью которого, можно создать платеж(пополнение счета)
Есть енпоинт, по которому создается платежь в платежной системе stripe
И после успешной оплаты - сумма, на которую произошла транзакиция, пополняется на баланс пользователя
Есть эндпоинты:
-Создание нового пользователя с нулевым балансом (OK) Отдельным эндпоинтом создания пользователя
-Проверка баланса пользователя по уникальному идентефикатору(username)
Просто гет запрос в который передается идентификатор пользователя
-Создание платежа в системе
-Проверка статуса транзакции пользователя


# PACKAGE:
gorm:  gorm.io/gorm
chi:    go get -u github.com/go-chi/chi/v5  
        go get -u github.com/go-chi/chi/v5/middleware  
        go get -u gorm.io/driver/sqlite
viper:  github.com/spf13/viper
logrus:  github.com/sirupsen/logrus
jwt(sha256):  go get -u github.com/golang-jwt/jwt/v5
godotenv:  github.com/joho/godotenv
pq: go get github.com/lib/pq
stripe: go get -u github.com/stripe/stripe-go/v74 
        go get -u github.com/stripe/stripe-go/v74/checkout/session
migration: go get -u github.com/golang-migrate/migrate/v4

Должна быть отдельная джоба которая входит в бд и подтягивает (статус)пендинг с пеймент сешн 
и проверяла состояние 



# Done
-Create base struct project +
-Server +
-Module + Setting output/balance +
-Connect to db and create migrations+-
-Dependency Injection(Внедрение зависимостей)+
-Implemented the endpoint "/userCreate" +
-Implemented the endpoint "/balance" +
-Connect stripe+
-Implemented the endpoint "/{username}/{sum}", h.createPaymentSession +
-Implemented the endpoint "/{sessionID}/status", h.getSessionStatus +


# Plans
1.
2.
3. 

# Тестовые карты:

Успешный платеж с картой:

Номер карты: 4242 4242 4242 4242
Срок действия: Любая будущая дата (например, 12/34)
CVC: 123
Платеж пройдет успешно.
Платеж с ошибкой (недостаточно средств):

Номер карты: 4000 0000 0000 9995
Срок действия: Любая будущая дата
CVC: 123
Ошибка: "Your card has insufficient funds."
Тестирование отказа по причине неправильного CVC:

Номер карты: 4000 0000 0000 0127
Срок действия: Любая будущая дата
CVC: 000
Ошибка: "Your card's security code is incorrect."
Тестирование отмены платежа:

Номер карты: 4000 0000 0000 0002
Срок действия: Любая будущая дата
CVC: 123
Ошибка: "Your card has been declined."
Тестирование карты с подтверждением 3D Secure:

Номер карты: 4000 0025 0000 3155
Срок действия: Любая будущая дата
CVC: 123
Поведение: Требуется 3D Secure (для тестирования 3D Secure).# stripe-project2





..
тестові параметри, для тесту інтеграції:
payee_id: 1185
Login: WDISHOP
Password: wdi451  (прохання не змінювати)
Посилання для входу в особистий кабінет: https://business.portmone.com.ua/login
Номер карти для проведення успішних оплат у тестовому режимі: 4444333322221111. Термін дії:
будь-яка дата не раніше поточного дня. CVV2-код: будь-який.
Для отримання помилки на стандартній сторінці оплати Portmone.com необхідно ввести такі реквізити платіжної картки:
4111111111111111.
Термін дії: будь-яка дата не раніше поточного дня.
CVV2-код: будь-який.