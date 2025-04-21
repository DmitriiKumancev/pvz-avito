# gRPC API для ПВЗ

## Описание

gRPC API для работы с ПВЗ (пунктами выдачи заказов). Позволяет получать список
ПВЗ через gRPC.

## Proto-файл

Для описания сервиса используется Protocol Buffers:

```protobuf
syntax = "proto3";

package pvz.v1;

option go_package = "github.com/dkumancev/avito-pvz/pkg/infrastructure/grpc/pb;pb";

import "google/protobuf/timestamp.proto";

service PVZService {
  rpc GetPVZList(GetPVZListRequest) returns (GetPVZListResponse);
}

message PVZ {
  string id = 1;
  google.protobuf.Timestamp registration_date = 2;
  string city = 3;
}

enum ReceptionStatus {
  RECEPTION_STATUS_IN_PROGRESS = 0;
  RECEPTION_STATUS_CLOSED = 1;
}

message GetPVZListRequest {}

message GetPVZListResponse {
  repeated PVZ pvzs = 1;
}
```

## Запуск gRPC сервера

Для запуска gRPC сервера используйте команду:

```bash
make run-grpc
```

По умолчанию сервер запускается на порту 50051.

## Пример клиента

В директории `examples/grpc_client` находится пример клиента, который
подключается к gRPC серверу и получает список ПВЗ.

Чтобы запустить клиент:

```bash
go run examples/grpc_client/main.go
```

## Тестирование с помощью grpcurl

Также можно использовать инструмент `grpcurl` для тестирования gRPC API:

```bash
# Список доступных сервисов
grpcurl -plaintext localhost:50051 list

# Список методов PVZ сервиса
grpcurl -plaintext localhost:50051 list pvz.v1.PVZService

# Получение списка ПВЗ
grpcurl -plaintext localhost:50051 pvz.v1.PVZService/GetPVZList
```
