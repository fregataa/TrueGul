# Sprint 3: Mobile App

## Overview

| 항목 | 내용 |
|------|------|
| 목표 | 핵심 UI 구현 (인증, 제출, 결과) |
| 선행 조건 | Sprint 0 완료, Sprint 1 API 스펙 확정 |
| 후속 Sprint | Sprint 4 (Integration) |

---

## Tasks

| ID | Task | 상태 | 비고 |
|----|------|------|------|
| S3-1 | 프로젝트 구조 설정 | TODO | 폴더 구조, 상태관리 |
| S3-2 | API Client 구현 | TODO | 인터셉터, 에러 처리 |
| S3-3 | Auth 상태 관리 | TODO | 토큰 저장, 자동 로그인 |
| S3-4 | Splash Screen | TODO | |
| S3-5 | Login/Signup Screen | TODO | |
| S3-6 | Home Screen | TODO | |
| S3-7 | Submit Screen | TODO | 글자수 카운터 |
| S3-8 | Submitted Screen | TODO | |
| S3-9 | History Screen | TODO | 페이지네이션 |
| S3-10 | Result Detail Screen | TODO | 점수 시각화 |
| S3-11 | Profile Screen | TODO | |
| S3-12 | Push Notification 연동 | TODO | FCM/APNs |
| S3-13 | Deep Link 처리 | TODO | 알림 → 결과 화면 |

---

## Screen Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Splash    │────▶│    Login    │────▶│    Home     │
└─────────────┘     └─────────────┘     └──────┬──────┘
                           │                   │
                           ▼                   │
                    ┌─────────────┐            │
                    │   Signup    │            │
                    └─────────────┘            │
                                               │
          ┌────────────────────────────────────┼────────────────────────────────────┐
          ▼                                    ▼                                    ▼
   ┌─────────────┐                      ┌─────────────┐                      ┌─────────────┐
   │   Submit    │                      │   History   │                      │   Profile   │
   │   (작성)    │                      │   (기록)    │                      │   (설정)    │
   └──────┬──────┘                      └──────┬──────┘                      └─────────────┘
          │                                    │
          ▼                                    ▼
   ┌─────────────┐                      ┌─────────────┐
   │  Submitted  │─────────────────────▶│   Result    │
   │  (제출완료) │     Push 알림 탭     │   Detail    │
   └─────────────┘                      └─────────────┘
```

---

## S3-1: Project Structure

### Flutter 구조

```
mobile/
├── lib/
│   ├── main.dart
│   ├── app.dart
│   │
│   ├── config/
│   │   ├── routes.dart
│   │   ├── theme.dart
│   │   └── constants.dart
│   │
│   ├── core/
│   │   ├── api/
│   │   │   ├── api_client.dart
│   │   │   ├── api_endpoints.dart
│   │   │   └── interceptors/
│   │   ├── errors/
│   │   │   └── app_exception.dart
│   │   └── storage/
│   │       └── secure_storage.dart
│   │
│   ├── features/
│   │   ├── auth/
│   │   │   ├── data/
│   │   │   ├── domain/
│   │   │   └── presentation/
│   │   ├── submission/
│   │   │   ├── data/
│   │   │   ├── domain/
│   │   │   └── presentation/
│   │   └── profile/
│   │       └── ...
│   │
│   ├── shared/
│   │   ├── widgets/
│   │   └── utils/
│   │
│   └── services/
│       ├── push_notification_service.dart
│       └── deep_link_service.dart
│
├── assets/
│   ├── images/
│   └── fonts/
│
├── android/
├── ios/
├── pubspec.yaml
└── README.md
```

### React Native 구조

```
mobile/
├── src/
│   ├── App.tsx
│   │
│   ├── config/
│   │   ├── navigation.tsx
│   │   ├── theme.ts
│   │   └── constants.ts
│   │
│   ├── api/
│   │   ├── client.ts
│   │   ├── endpoints.ts
│   │   └── interceptors.ts
│   │
│   ├── features/
│   │   ├── auth/
│   │   │   ├── screens/
│   │   │   ├── hooks/
│   │   │   └── types.ts
│   │   ├── submission/
│   │   └── profile/
│   │
│   ├── components/
│   │   └── common/
│   │
│   ├── hooks/
│   ├── stores/
│   ├── utils/
│   └── services/
│       ├── pushNotification.ts
│       └── deepLink.ts
│
├── android/
├── ios/
├── package.json
└── README.md
```

---

## S3-2: API Client

### Endpoints

| Method | Endpoint | 용도 |
|--------|----------|------|
| POST | `/api/v1/auth/signup` | 회원가입 |
| POST | `/api/v1/auth/login` | 로그인 |
| POST | `/api/v1/auth/logout` | 로그아웃 |
| GET | `/api/v1/auth/me` | 현재 사용자 |
| POST | `/api/v1/submissions` | 답안 제출 |
| GET | `/api/v1/submissions/:id` | 결과 조회 |
| GET | `/api/v1/submissions` | 제출 이력 |
| POST | `/api/v1/push/register` | Push 토큰 등록 |

### API Client (Flutter/Dio)

```dart
// lib/core/api/api_client.dart
import 'package:dio/dio.dart';

class ApiClient {
  late final Dio _dio;
  final SecureStorage _storage;

  ApiClient(this._storage) {
    _dio = Dio(BaseOptions(
      baseUrl: AppConstants.apiBaseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 30),
    ));

    _dio.interceptors.addAll([
      AuthInterceptor(_storage),
      LogInterceptor(),
      ErrorInterceptor(),
    ]);
  }

  Future<Response<T>> get<T>(String path, {Map<String, dynamic>? params}) {
    return _dio.get(path, queryParameters: params);
  }

  Future<Response<T>> post<T>(String path, {dynamic data}) {
    return _dio.post(path, data: data);
  }

  Future<Response<T>> delete<T>(String path) {
    return _dio.delete(path);
  }
}

// Auth Interceptor
class AuthInterceptor extends Interceptor {
  final SecureStorage _storage;

  AuthInterceptor(this._storage);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) async {
    final token = await _storage.getToken();
    if (token != null) {
      options.headers['Cookie'] = 'token=$token';
    }
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    // Extract and save token from Set-Cookie header
    final cookies = response.headers['set-cookie'];
    if (cookies != null) {
      for (final cookie in cookies) {
        if (cookie.startsWith('token=')) {
          final token = cookie.split(';')[0].substring(6);
          _storage.saveToken(token);
        }
      }
    }
    handler.next(response);
  }
}
```

---

## S3-3: Auth State Management

### Flutter (Riverpod)

```dart
// lib/features/auth/domain/auth_state.dart
import 'package:freezed_annotation/freezed_annotation.dart';

part 'auth_state.freezed.dart';

@freezed
class AuthState with _$AuthState {
  const factory AuthState.initial() = _Initial;
  const factory AuthState.loading() = _Loading;
  const factory AuthState.authenticated(User user) = _Authenticated;
  const factory AuthState.unauthenticated() = _Unauthenticated;
  const factory AuthState.error(String message) = _Error;
}

// lib/features/auth/domain/auth_notifier.dart
class AuthNotifier extends StateNotifier<AuthState> {
  final AuthRepository _repository;
  final SecureStorage _storage;

  AuthNotifier(this._repository, this._storage) : super(const AuthState.initial());

  Future<void> checkAuthStatus() async {
    state = const AuthState.loading();
    try {
      final user = await _repository.getCurrentUser();
      state = AuthState.authenticated(user);
    } catch (e) {
      state = const AuthState.unauthenticated();
    }
  }

  Future<void> login(String email, String password) async {
    state = const AuthState.loading();
    try {
      final user = await _repository.login(email, password);
      state = AuthState.authenticated(user);
    } catch (e) {
      state = AuthState.error(e.toString());
    }
  }

  Future<void> logout() async {
    await _repository.logout();
    await _storage.clearToken();
    state = const AuthState.unauthenticated();
  }
}
```

---

## S3-4 ~ S3-11: Screen Specifications

### Splash Screen

| 항목 | 상세 |
|------|------|
| 기능 | 앱 로딩, 자동 로그인 체크 |
| 이동 | 로그인 됨 → Home, 아님 → Login |
| 디자인 | 로고 중앙, 로딩 인디케이터 |

### Login Screen

| 항목 | 상세 |
|------|------|
| 입력 | 이메일, 비밀번호 |
| 버튼 | 로그인, 회원가입 이동 |
| 유효성 | 이메일 형식, 비밀번호 최소 6자 |
| 에러 | 인라인 에러 메시지 |

### Signup Screen

| 항목 | 상세 |
|------|------|
| 입력 | 이메일, 비밀번호, 비밀번호 확인 |
| 버튼 | 회원가입, 로그인 이동 |
| 유효성 | 이메일 형식, 비밀번호 일치 |

### Home Screen

| 항목 | 상세 |
|------|------|
| 섹션 | 새 채점 시작 CTA, 최근 기록 3개 |
| 하단 탭 | 홈, 기록, 프로필 |
| 액션 | 새 채점 → Submit, 기록 항목 → Result Detail |

### Submit Screen

```
┌─────────────────────────────────────┐
│ ← 취소            TOPIK II 54번    │
├─────────────────────────────────────┤
│                                     │
│  [문제]                             │
│  다음을 주제로 하여 자신의 생각을   │
│  600~700자로 쓰십시오.              │
│                                     │
│  "현대 사회에서 환경 보호의         │
│   중요성에 대해 논하시오."          │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  ┌─────────────────────────────┐   │
│  │                             │   │
│  │  [텍스트 입력 영역]         │   │
│  │                             │   │
│  │                             │   │
│  │                             │   │
│  └─────────────────────────────┘   │
│                                     │
│  글자 수: 0 / 700                   │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  [        채점 요청        ]        │
│                                     │
└─────────────────────────────────────┘
```

| 항목 | 상세 |
|------|------|
| 입력 | 다중 라인 텍스트 (최대 800자) |
| 실시간 | 글자 수 카운터 |
| 유효성 | 100자 이상 필수, 800자 이하 |
| 버튼 | 채점 요청 (유효할 때만 활성화) |

### Submitted Screen

```
┌─────────────────────────────────────┐
│              채점 중                │
├─────────────────────────────────────┤
│                                     │
│         [로딩 애니메이션]           │
│                                     │
│      AI가 답안을 분석하고           │
│      있습니다...                    │
│                                     │
│      예상 소요 시간: 약 30초        │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  ℹ️ 앱을 닫아도 채점이 진행됩니다.  │
│     완료되면 알림으로 알려드려요.   │
│                                     │
├─────────────────────────────────────┤
│                                     │
│  [      채점 기록으로 이동      ]   │
│                                     │
└─────────────────────────────────────┘
```

| 항목 | 상세 |
|------|------|
| 표시 | 로딩 애니메이션, 안내 메시지 |
| 백그라운드 | 앱 종료 가능 안내 |
| 버튼 | 채점 기록으로 이동 |

### History Screen

```
┌─────────────────────────────────────┐
│             채점 기록               │
├─────────────────────────────────────┤
│ [전체] [대기중] [완료]              │
├─────────────────────────────────────┤
│ ┌─────────────────────────────────┐ │
│ │ NEW  2026.01.09 14:30           │ │
│ │ 총점: 42점                      │ │
│ │ 현대 사회에서 환경 보호의...    │ │
│ └─────────────────────────────────┘ │
│ ┌─────────────────────────────────┐ │
│ │      2026.01.08 10:15           │ │
│ │ 총점: 38점                      │ │
│ │ 기술 발전이 인간 관계에...      │ │
│ └─────────────────────────────────┘ │
│ ┌─────────────────────────────────┐ │
│ │ ⏳   2026.01.09 15:00           │ │
│ │ 채점 중...                      │ │
│ │ 독서의 중요성에 대해...         │ │
│ └─────────────────────────────────┘ │
│                                     │
│         [더 보기]                   │
│                                     │
└─────────────────────────────────────┘
```

| 항목 | 상세 |
|------|------|
| 필터 | 전체, 대기중, 완료 |
| 리스트 | 날짜, 점수, 답안 미리보기 |
| 배지 | NEW (미확인), ⏳ (대기중) |
| 페이지네이션 | 무한 스크롤 또는 더보기 버튼 |

### Result Detail Screen

```
┌─────────────────────────────────────┐
│ ←                 채점 결과         │
├─────────────────────────────────────┤
│                                     │
│         총점: 42 / 50               │
│         추정 등급: 5급              │
│                                     │
│  ┌─────────────────────────────┐   │
│  │   [점수 차트/게이지]         │   │
│  │   내용: 16/20                │   │
│  │   구성: 13/15                │   │
│  │   표현: 13/15                │   │
│  └─────────────────────────────┘   │
│                                     │
│  ⚠️ AI 작성 의심 (선택적 표시)      │
│                                     │
├─────────────────────────────────────┤
│  📝 내용 및 과제 수행               │
│  ─────────────────────────────────  │
│  주제에 맞게 자신의 의견을 명확히   │
│  제시했습니다. 다만 구체적인 예시가 │
│  부족합니다...                      │
│                                     │
│  📐 글의 전개 구조                  │
│  ─────────────────────────────────  │
│  서론-본론-결론 구조가 잘 갖춰져    │
│  있습니다...                        │
│                                     │
│  💬 언어 사용                       │
│  ─────────────────────────────────  │
│  중급 수준의 어휘를 적절히 사용...  │
│                                     │
│  📋 종합 피드백                     │
│  ─────────────────────────────────  │
│  전반적으로 TOPIK II 5급 수준의     │
│  글입니다...                        │
│                                     │
├─────────────────────────────────────┤
│  [원본 답안 보기]                   │
└─────────────────────────────────────┘
```

| 항목 | 상세 |
|------|------|
| 헤더 | 총점, 추정 등급 |
| 차트 | 항목별 점수 시각화 (게이지/바) |
| 피드백 | 4개 영역 (내용, 구성, 표현, 종합) |
| AI 경고 | 플래그 시에만 표시 |
| 액션 | 원본 답안 보기 (모달) |

### Profile Screen

| 항목 | 상세 |
|------|------|
| 정보 | 이메일 표시 |
| 설정 | 알림 설정 (On/Off) |
| 액션 | 로그아웃 버튼 |

---

## S3-12: Push Notification

### Flutter (firebase_messaging)

```dart
// lib/services/push_notification_service.dart
import 'package:firebase_messaging/firebase_messaging.dart';

class PushNotificationService {
  final FirebaseMessaging _messaging = FirebaseMessaging.instance;
  final ApiClient _api;

  PushNotificationService(this._api);

  Future<void> initialize() async {
    // Request permission
    final settings = await _messaging.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    if (settings.authorizationStatus == AuthorizationStatus.authorized) {
      // Get FCM token
      final token = await _messaging.getToken();
      if (token != null) {
        await _registerToken(token);
      }

      // Listen for token refresh
      _messaging.onTokenRefresh.listen(_registerToken);

      // Handle foreground messages
      FirebaseMessaging.onMessage.listen(_handleForegroundMessage);

      // Handle background/terminated messages
      FirebaseMessaging.onMessageOpenedApp.listen(_handleMessageOpen);
    }
  }

  Future<void> _registerToken(String token) async {
    await _api.post('/api/v1/push/register', data: {
      'token': token,
      'platform': Platform.isIOS ? 'ios' : 'android',
    });
  }

  void _handleForegroundMessage(RemoteMessage message) {
    // Show local notification
    if (message.data['type'] == 'scoring_complete') {
      _showLocalNotification(
        title: message.notification?.title ?? '채점 완료',
        body: message.notification?.body ?? '결과를 확인하세요',
        payload: message.data['submission_id'],
      );
    }
  }

  void _handleMessageOpen(RemoteMessage message) {
    // Navigate to result detail
    if (message.data['type'] == 'scoring_complete') {
      final submissionId = message.data['submission_id'];
      NavigationService.navigateTo('/result/$submissionId');
    }
  }
}
```

---

## S3-13: Deep Link

```dart
// lib/services/deep_link_service.dart
class DeepLinkService {
  void initialize() {
    // Handle initial link (app opened via deep link)
    getInitialLink().then(_handleDeepLink);

    // Listen for incoming links while app is running
    linkStream.listen(_handleDeepLink);
  }

  void _handleDeepLink(String? link) {
    if (link == null) return;

    final uri = Uri.parse(link);

    // truegul://result/{submission_id}
    if (uri.host == 'result' && uri.pathSegments.isNotEmpty) {
      final submissionId = uri.pathSegments.first;
      NavigationService.navigateTo('/result/$submissionId');
    }
  }
}
```

---

## Completion Criteria

- [ ] 프로젝트 구조 설정 완료
- [ ] API Client 구현 및 인터셉터 설정
- [ ] Auth 상태 관리 구현
- [ ] 모든 화면 UI 구현
  - [ ] Splash
  - [ ] Login/Signup
  - [ ] Home
  - [ ] Submit
  - [ ] Submitted
  - [ ] History
  - [ ] Result Detail
  - [ ] Profile
- [ ] Push Notification 연동
- [ ] Deep Link 처리
- [ ] 기본 에러 처리 (네트워크, 인증)
- [ ] 로딩 상태 처리

---

*Sprint 1 API 완료 후 실제 연동 테스트 진행*
