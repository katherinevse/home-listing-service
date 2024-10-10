# Сервис домов

Ежедневно в нашем сервисе недвижимости публикуются тысячи объявлений о продаже или аренде. 
Пользователи могут выбирать жильё по нужным параметрам в понравившемся доме. 
Прежде чем попасть в каталог, каждое объявление проходит тщательную модерацию.

## Описание проекта

В этом проекте я разработала бэкенд-сервис, который позволяет пользователям продавать квартиры, загружая объявления.
Сервис включает в себя функционал авторизации, создания домов и квартир, модерации объявлений,
а также возможность подписки на уведомления о новых квартирах.

**Kafka** используется для асинхронной обработки уведомлений, что позволяет пользователям получать обновления о новых объявлениях в реальном времени.

## Для реализации данного сервиса я создала Miro и kanban доски, которыми пользовалась в момент разработки :
**Ссылка на Miro и kanban для разработки: https://miro.com/app/board/uXjVLe9m1Lc=/?share_link_id=432262272261**

## Функционал сервиса:

1. **Авторизация пользователей**:
	- Используйте [endpoint /dummyLogin](https://github.com/katherinevse/home-listing-service/blob/main/api.yaml#L10) для получения токена с уровнем доступа:
	  ```bash
      curl -X POST "http://localhost:8080/dummyLogin?userType=client"
      ```

	- Регистрация и авторизация пользователей по почте и паролю:
		- Регистрация:
		  ```bash
          curl -X POST "http://localhost:8080/register" -d '{"email":"example@mail.com", "password":"yourpassword", "userType":"client"}'
          ```
		- Авторизация:
		  ```bash
          curl -X POST "http://localhost:8080/login" -d '{"email":"example@mail.com", "password":"yourpassword"}'
          ```

2. **Создание дома**:
	- Только модератор может создать дом через [endpoint /house/create](https://github.com/katherinevse/home-listing-service/blob/main/api.yaml#L20):
	  ```bash
      curl -X POST "http://localhost:8080/house/create" -H "Authorization: Bearer <your_token>" -d '{"address":"123 Main St", "yearBuilt":2000}'
      ```

3. **Создание квартиры**:
	- Любой пользователь может создать квартиру через [endpoint /flat/create](https://github.com/katherinevse/home-listing-service/blob/main/api.yaml#L30):
	  ```bash
      curl -X POST "http://localhost:8080/flat/create" -H "Authorization: Bearer <your_token>" -d '{"houseID":1, "flatNumber":101, "price":50000, "rooms":2}'
      ```

4. **Модерация квартиры**:
	- Модератор может изменить статус квартиры через [endpoint /flat/update](https://github.com/katherinevse/home-listing-service/blob/main/api.yaml#L40):
	  ```bash
      curl -X PUT "http://localhost:8080/flat/update" -H "Authorization: Bearer <moderator_token>" -d '{"flatID":1, "status":"approved"}'
      ```

5. **Получение списка квартир по номеру дома**:
	- Получение списка квартир через [endpoint /house/{id}](https://github.com/katherinevse/home-listing-service/blob/main/api.yaml#L50):
	  ```bash
      curl -X GET "http://localhost:8080/house/1" -H "Authorization: Bearer <your_token>"
      ```

6. **Подписка на уведомления**:
	- Подписка на уведомления о новых квартирах через [endpoint /house/{id}/subscribe](https://github.com/katherinevse/home-listing-service/blob/main/api.yaml#L60):
	  ```bash
      curl -X POST "http://localhost:8080/house/1/subscribe" -H "Authorization: Bearer <your_token>"
      ```

## Общие вводные

- **Дом**:
	- Уникальный номер
	- Адрес
	- Год постройки
	- Застройщик
	- Дата создания
	- Дата последнего добавления

- **Квартира**:
	- Номер квартиры
	- Цена
	- Количество комнат

## Условия

Проект полностью соответствует [API спецификации](https://github.com/katherinevse/home-listing-service/blob/main/api.yaml) и включает в себя интеграционные и модульные тесты для сценариев получения списка квартир и публикации новой квартиры.

## Дополнительные задания


## Дополнительно

- Реализована пользовательская авторизация по методам /register и /login.
- Настроен асинхронный механизм уведомления пользователя о появлении новых квартир с использованием Kafka. 
- Интеграция middleware обеспечивает проверку авторизации пользователей при доступе к защищенным ресурсам, что улучшает безопасность сервиса.
- Использован логгер  **slog** с разлияными уровнями для ведения журналов ошибок и событий, что позволяет отслеживать действия в системе и улучшает управление логами.

[//]: # (- Настроены CI и кодогенерация DTO endpoint'ов по openapi схеме.)


## Запуск
Сервис можно запустить с помощью Docker. В корне репозитория выполните:
```bash
docker-compose up --build
