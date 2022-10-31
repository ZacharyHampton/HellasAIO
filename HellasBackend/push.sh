docker buildx build --platform=linux/amd64 -t hellasbackend .
aws ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin aws_repository
docker tag hellasbackend:latest aws_repository:latest
docker push aws_repository
aws ecs update-service --cluster backend-api --service HellasService --region eu-central-1 --force-new-deployment