# ASCII Player

ASCII 아트 애니메이션을 재생하는 간단한 Go 프로그램입니다.

## 기능

- 텍스트 파일에서 ASCII 프레임 로드
- 설정 가능한 프레임 레이트로 애니메이션 재생
- 간단한 명령줄 인터페이스

## 사용법

```bash
# 프로그램 빌드
go build -o ascii-player

# 애니메이션 재생
./ascii-player examples/simple.txt
```

## 파일 형식

애니메이션 파일은 다음과 같은 형식을 따라야 합니다:

```
프레임 1 내용
---
프레임 2 내용
---
프레임 3 내용
```

각 프레임은 `---` 구분자로 분리됩니다.

## 프로젝트 구조

```
├── main.go              # 메인 진입점
├── internal/            # 내부 패키지
│   └── player.go        # 애니메이션 플레이어 로직
├── pkg/                 # 공개 패키지
│   └── loader.go        # 파일 로더 유틸리티
└── examples/            # 예제 애니메이션 파일
    └── simple.txt       # 간단한 애니메이션 예제
```

## 개발

```bash
# 의존성 설치
go mod tidy

# 테스트 실행
go test ./...

# 프로그램 실행
go run main.go examples/simple.txt
```
