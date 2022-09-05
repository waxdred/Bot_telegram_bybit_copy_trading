
NAME = bot
MOD = go.mod
SUM = go.sum


all: $(NAME)

$(NAME): build
	./${NAME}

install: go.mod go.sum 
	@-go mod init bot
	@-go mod tidy

build: 
	@-go build .

fclean:
	rm ./${MOD}
	rm ./${SUM}
	rm ./${NAME}

re: fclean all
