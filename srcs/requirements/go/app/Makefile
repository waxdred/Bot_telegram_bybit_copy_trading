
NAME = bot
MOD = go.mod
SUM = go.sum
all: $(NAME)

$(NAME): build
	./${NAME}

build: 
	@-go mod init bot
	@-go mod tidy
	@-go build .

fclean:
	rm ./${MOD}
	rm ./${SUM}
	rm ./${NAME}

.PHONY: all fclean
.SILENT :
re: fclean all
