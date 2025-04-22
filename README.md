## Запуск 
В докере: docker-compose up  
Raw: go run ./main migrate && go run ./main serve

## API
https://www.postman.com/red-water-310799/todo/overview  
Там есть примеры всех рутов + оттуда можно потестить, считаю, что это сойдет за документацию.

## Проблемы
Почему-то в контейреы логи печатаются в os.Stdou по 2 раза, почему разбираться времени не было
