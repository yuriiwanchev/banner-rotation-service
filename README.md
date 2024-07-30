# Cервис "Ротация баннеров"

![example workflow](https://github.com/yuriiwanchev/banner-rotation-service/actions/workflows/linter.yaml/badge.svg)
![example workflow](https://github.com/yuriiwanchev/banner-rotation-service/actions/workflows/tests.yaml/badge.svg)
![example workflow](https://github.com/yuriiwanchev/banner-rotation-service/actions/workflows/build.yaml/badge.svg)

## Общее описание

Сервис "Ротация баннеров" предназначен для выбора наиболее эффективных (кликабельных) баннеров, в условиях меняющихся предпочтений пользователей и набора баннеров.

Задача сервиса - осуществлять "ротацию" баннеров, показывая те, которые наиболее вероятно приведут к переходу. Для этого используется алгоритм "Многорукий бандит".

Сервис предоставляет REST API.

Также микросервис отправляет события кликов и показов в брокер сообщений Kafka для дальнейшей обработки в аналитических системах.

## Развертывание сервиса

Развертывание микросервиса должно осуществляться командой make run в директории с проектом (banner-rotation-service).

Проект возможно собрать через make build

## Тестирование

Для юнит тестов можно использовать `make test` в директории с проектом.

Интеграционные тесты, проверяющие работу сервиса через его API, запускаются через команду `make integration-test`.
