all: start


start-docker:
	@docker compose up -d

stop-docker:
	@docker compose down

deploy-kuber:
	@kubectl apply -f main-namespace.yaml
	@kubectl -n main apply -f ingress.yaml
	@kubectl -n main apply -f user-service/user-deployment.yaml
	@kubectl -n main apply -f user-service/user-service.yaml

	@kubectl -n main apply -f auth-service/auth-deployment.yaml
	@kubectl -n main apply -f auth-service/auth-service.yaml

	@kubectl -n main apply -f messaging-service/messaging-deployment.yaml
	@kubectl -n main apply -f messaging-service/messaging-service.yaml

	@kubectl -n main apply -f notification-service/notification-deployment.yaml
	@kubectl -n main apply -f notification-service/notification-service.yaml

	@kubectl -n main apply -f server-service/server-deployment.yaml
	@kubectl -n main apply -f server-service/server-service.yaml


