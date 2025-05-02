# Quiz Web Application Development Plan

This document outlines the tasks required to build the Vue.js frontend and Golang backend quiz application.

## Phase 1: Project Setup & Backend Foundation

- [x] Create project directory: `quiz-app`
  - [x] Create `backend` subdirectory
  - [x] Create `frontend` subdirectory
- [x] Initialize Go backend (`backend`)
  - [x] `go mod init <module_path>`
  - [x] Choose web framework (e.g., Gin) and add dependency
  - [x] Choose database & ORM (e.g., SQLite + GORM) and add dependencies
  - [x] Basic `main.go` setup
  - [x] Setup environment variable loading (e.g., godotenv)
- [x] Initialize Vue.js frontend (`frontend`)
  - [x] Use Vite: `npm create vite@latest frontend -- --template vue-ts` (or Vue CLI)
  - [x] Install Vue Router, Pinia, Axios
  - [x] Optional: Install UI Framework (UnoCSS instead of Tailwind)
- [ ] Define Database Schema
  - [ ] `admin_users` (id, username, password_hash)
  - [ ] `quizzes` (id, title, description, time_limit_seconds, created_by, created_at, updated_at)
  - [ ] `questions` (id, quiz_id, text, type [single/multi], order, created_at, updated_at)
  - [ ] `choices` (id, question_id, text, is_correct, order, created_at, updated_at)
  - [ ] `responder_credentials` (id, quiz_id, username, password_hash, expires_at, created_at)
  - [ ] `responses` (id, quiz_id, responder_username, submitted_at, score)
  - [ ] `answers` (id, response_id, question_id, selected_choice_ids_json)
- [x] Implement Backend Models (Go structs in `models/`)
- [x] Implement Database Layer (`database/`)
  - [x] Connection setup
  - [x] Auto-migration function
  - [ ] Basic CRUD functions (or rely on GORM directly initially)
- [x] Setup Basic API Routing (Gin/Echo in `main.go` or `routes/`)
- [x] Implement Admin Authentication (`auth/`)
  - [x] Password hashing utility
  - [x] JWT generation utility
  - [x] Admin Login endpoint (`POST /admin/login`)
  - [x] JWT Auth middleware
- [ ] Implement Temporary Responder Authentication (`auth/`)
  - [x] Credential generation utility (random string, hash, expiry)
  - [x] Responder Login endpoint (`POST /responder/login`)
  - [ ] Middleware/check for responder routes (validate credentials, check expiry)
- [x] Configure CORS middleware
- [ ] Create initial `README.md`

## Phase 2: Core Admin Functionality (Backend API)

- [x] API Endpoint: Create Quiz (`POST /admin/quizzes`)
  - [x] Handler logic (receive quiz data, questions, choices; save to DB)
  - [x] Associate with logged-in admin
- [x] API Endpoint: Get Quizzes (`GET /admin/quizzes`)
  - [x] Handler logic (fetch quizzes created by the admin)
- [x] API Endpoint: Get Quiz Details (`GET /admin/quizzes/{quiz_id}`)
  - [x] Handler logic (fetch specific quiz with questions and choices, including correct answers)
- [x] API Endpoint: Update Quiz (`PUT /admin/quizzes/{quiz_id}`)
  - [x] Handler logic (update quiz title, description, time limit)
- [x] API Endpoint: Delete Quiz (`DELETE /admin/quizzes/{quiz_id}`)
  - [x] Handler logic (delete quiz and associated questions/choices/credentials/responses)
- [x] API Endpoint: Add Question to Quiz (`POST /admin/quizzes/{quiz_id}/questions`)
  - [x] Handler logic (add a new question with choices to an existing quiz)
- [x] API Endpoint: Update Question (`PUT /admin/questions/{question_id}`)
  - [x] Handler logic (update question text, type, order, choices)
- [x] API Endpoint: Delete Question (`DELETE /admin/questions/{question_id}`)
  - [x] Handler logic (delete a specific question and its choices)
- [x] API Endpoint: Generate Temporary Credentials (`POST /admin/quizzes/{quiz_id}/credentials`)
  - [x] Handler logic (generate username/password, set expiry, store hash)
- [x] API Endpoint: View Credentials (`GET /admin/quizzes/{quiz_id}/credentials`)
  - [x] Handler logic (list active credentials for a quiz)
- [x] API Endpoint: Revoke Credential (`DELETE /admin/credentials/{credential_id}`)
  - [x] Handler logic (delete a specific credential)
- [x] API Endpoint: View Responses (`GET /admin/quizzes/{quiz_id}/responses`)
  - [x] Handler logic (list all responses submitted for a quiz, including responder username and score)
- [x] API Endpoint: View Specific Response Details (`GET /admin/responses/{response_id}`)
  - [x] Handler logic (fetch details of a single response, including answers given)

## Phase 3: Core Responder Functionality (Backend API)

- [ ] API Endpoint: Get Quiz for Responder (`GET /quizzes/{quiz_id}`)
  - [ ] Requires temporary credential auth (use middleware/check)
  - [ ] Handler logic (fetch quiz details, questions, choices _without_ `is_correct` flag)
- [ ] API Endpoint: Submit Quiz (`POST /quizzes/{quiz_id}/submit`)
  - [ ] Requires temporary credential auth
  - [ ] Handler logic:
    - [ ] Receive answers
    - [ ] Validate submission (e.g., within time limit if applicable - maybe check server-side start time? Or trust client timer initially?)
    - [ ] Calculate score based on submitted answers vs correct choices
    - [ ] Store response record (username, quiz_id, score, timestamp)
    - [ ] Store answers given
    - [ ] Return score and correct answers map (question_id -> list of correct choice_ids)

## Phase 4: Frontend - Admin UI

- [x] Setup Vue Router for Admin section (`/admin/*`)
- [x] Setup Pinia stores (`authStore`, `adminQuizStore`)
- [x] Create API client service (`src/services/api.js`) with admin functions
- [x] Implement Admin Login Page (`AdminLogin.vue`)
  - [x] Form, call login API, store token, redirect
- [x] Implement Admin Layout/Dashboard (`AdminDashboard.vue`)
  - [x] Navigation (Quiz List, Create Quiz etc.)
  - [x] Logout button
  - [x] Router outlet
- [x] Implement Quiz List View (`AdminQuizList.vue`)
  - [x] Fetch and display admin's quizzes
  - [x] Links/buttons for Edit, Delete, View Responses, Manage Credentials (placeholders exist)
- [x] Implement Quiz Create/Edit Form (Component: `QuizForm.vue`, Create View: `AdminQuizCreate.vue`)
  - [x] Form for quiz details (title, description, time limit)
  - [x] Component for adding/editing questions (text, type [single/multi])
  - [x] Component for adding/editing choices (text, is_correct checkbox)
  - [x] Dynamic add/remove buttons for questions/choices
  - [x] Implement Quiz Edit View (`AdminQuizEdit.vue`) - Implemented
  - [x] API calls for create (in `AdminQuizCreate.vue`)
  - [x] API calls for update (in `AdminQuizEdit.vue`)
- [ ] Implement Credential Management View (`AdminCredentialManager.vue`) - *Needs creation*
  - [ ] Fetch/display existing credentials for a quiz
  - [ ] Form to generate new credentials (set expiry)
  - [ ] Display newly generated credentials
- [ ] Implement Response List View (`AdminResponseList.vue`) - *Needs creation*
  - [ ] Fetch/display responses for a quiz (responder, score, submitted_at)
  - [ ] Link to view response details
- [ ] Implement Response Details View (`AdminResponseDetails.vue`) - *Needs creation*
  - [ ] Fetch/display specific response details (answers given vs correct answers)

## Phase 5: Frontend - Responder UI

- [ ] Setup Vue Router for Responder section (`/login`, `/quiz/:id`, `/quiz/:id/result`)
- [ ] Update Pinia `authStore` to handle responder login state
- [ ] Update API client service with responder functions
- [ ] Implement Responder Login Page (`ResponderLogin.vue`)
  - [ ] Form for username/password, call login API, store temp token/session, redirect to quiz
- [ ] Implement Quiz Taking View (`QuizTaker.vue`)
  - [ ] Fetch quiz data based on route param `:id` (requires auth)
  - [ ] Display quiz title, description
  - [ ] Implement `Timer.vue` component and display if `time_limit_seconds > 0`
  - [ ] Implement `QuestionDisplay.vue` component
  - [ ] Render questions and choices using `QuestionDisplay.vue`
  - [ ] Store selected answers in local state
  - [ ] Handle submission: send answers to API, navigate to result page
  - [ ] Handle time expiry (auto-submit or notify user)
- [ ] Implement Quiz Result View (`QuizResult.vue`)
  - [ ] Get score and correct answers from route params or store after submission
  - [ ] Display score
  - [ ] Display questions again, highlighting user's answers and correct answers

## Phase 6: Refinement & Deployment

- [ ] Add Input Validation (Backend: handler level; Frontend: form level)
- [ ] Implement Comprehensive Error Handling (Backend: proper HTTP status codes, error responses; Frontend: user notifications)
- [ ] Styling and UI/UX Improvements (Use UI framework consistently, ensure responsiveness)
- [ ] Write Unit Tests (Backend: handlers, auth, logic)
- [ ] Write Component Tests (Frontend: key components like `QuizForm`, `QuestionDisplay`)
- [ ] Update `README.md` with final setup, run, and build instructions
- [ ] Configure Dockerfile for Backend (optional)
- [ ] Configure `docker-compose.yml` for local dev (optional)
- [ ] Setup CI/CD Pipeline (e.g., GitHub Actions)
  - [ ] Backend tests & build
  - [ ] Frontend tests & build
  - [ ] Deployment steps
- [ ] Deploy Backend (Cloud Run, ECS, Heroku, etc.)
- [ ] Deploy Frontend (Netlify, Vercel, S3/CloudFront, etc.)
- [ ] Configure Production Environment Variables (DB connection, JWT secret, frontend API URL)
