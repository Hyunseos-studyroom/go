# 미래의 나를 위한 카카오 oauth 정리

# For Web
## workflow
1. (프론트) 카카오 서버 Redirect URI로 요청을 보내서 인가코드 받기
2. (프론트) 받은 인가코드를 백엔드로 전달
3. (백엔드) 프론트에서 받은 인가코드를 사용해 카카오 서버로 토큰 요청
4. 카카오 서버에서 redirect uri, clientid, 인가코드를 검증 후 회원 정보가 담긴 카카오 토큰을 백엔드로 전송
5. (백엔드) 받은 토큰으로 유저정보를 받아서 우리 서버에 유저를 등록함
6. (백엔드) 카카오에서 제공하는 refresh 토큰이 아닌, 우리 서버의 JWT를 생성하여 프론트로 전송
7. (프론트) JWT 받은걸 쿠키에 저장해서 로그인