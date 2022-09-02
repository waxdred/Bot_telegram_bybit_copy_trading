NAME = bot

all:	$(NAME)

$(NAME):
	docker compose -v -f ./srcs/docker-compose.yml up --force-recreate --build -d 

down:
	docker compose -f ./srcs/docker-compose.yml down

init:
	@cp ./srcs/.env ./srcs/requirements/python/.
	@rm -rf ./srcs/requirements/python/config
	@-mkdir ./srcs/requirements/python/config
	@cd ./srcs/requirements/python/config && python3.10 ../app/telegram.py init
	@rm ./srcs/requirements/python/.env

fclean:
	@-docker rmi go
	@-docker system prune -a --force

re: down fclean all
