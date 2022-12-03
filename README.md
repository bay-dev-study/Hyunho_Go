## 진행 상태
[노마드 코인 강좌](https://nomadcoders.co/nomadcoin)에서 구현한 내용 모두 정상 작동하는 것까지 확인하였음.

test code 작성하는 부분 중 utils package만 coverage를 모두 채웠고, 그 이외의 package에 대한 test code는 작성하지 않고 여기서 스터디를 마칠 예정

## 실행 방법
go run main.go launch rest --port 포트번호

vscode 확장 프로그램 중 [rest api client](https://marketplace.visualstudio.com/items?itemName=humao.rest-client)를 설치한 뒤, 아래 내용을 복사 붙여넣기하여 사용하면 됨 (폴더 내 rest_api_check.http 참조)

위 방법 이외에도 postman / curl 등등 기타 http request를 보낼 수 있는 어떤 프로그램으로도 사용 가능

## RestApi 설명

localhost는 동일한 컴퓨터에서 여러개의 노드를 실행시키는 상황에서 입력하는 것이며, 
서로 다른 컴퓨터에서 노드를 실행시킬 때는 해당 컴퓨터의 외부ip를 localhost 대신에 입력해주면 됨

```
# rest api의 document 확인
GET http://localhost:포트번호/ HTTP/1.1   

# 해당 노드의 블록체인 현재 상태 확인
GET http://localhost:포트번호/status HTTP/1.1

# 해당 노드의 현재 블록 전부 보기
GET http://localhost:포트번호/blocks HTTP/1.1

# 해당 노드에 트랜잭션 등록 (서명을 위해서는 자기자신의 wallet 파일(key.wallet)이 동일 폴더 내에 있어야함)
POST http://localhost:포트번호/transactions HTTP/1.1
content-type: application/json

{
    "to": "월렛주소",
    "amount": 50
}

# 해당 노드의 mempool 보기
GET http://localhost:포트번호/mempool HTTP/1.1

# 해당 월렛 주소의 unspent transaction output 보기
GET http://localhost:포트번호/balance/월렛주소 HTTP/1.1

# 해당 월렛 주소의 잔고 보기
GET http://localhost:포트번호/balance/월렛주소?total=true HTTP/1.1

# 해당 노드 채굴 시작 
GET http://localhost:포트번호/confirm HTTP/1.1

# 해당 해시의 블록 확인
GET http://localhost:포트번호/blocks/해시 HTTP/1.1

# 해당 노드에 연결된 peer 확인
GET http://localhost:포트번호/peer HTTP/1.1

# 해당 노드에 peer 추가 
POST http://localhost:포트번호/peer HTTP/1.1
content-type: application/json

{
    "address": "127.0.0.1", 
    "port": "연결을 시키고자 하는 peer의 포트번호"
}
```
