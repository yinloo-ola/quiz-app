import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

// Import views/components here (we'll create them later)
// Example: import HomeView from '../views/HomeView.vue';
// Example: import AdminLogin from '../views/admin/AdminLogin.vue';
// Example: import AdminLayout from '../layouts/AdminLayout.vue';

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'home',
    // component: HomeView // We'll create this later
    component: () => import('../views/HomeView.vue'), // Lazy load home
  },
  {
    path: '/admin/login',
    name: 'admin-login',
    // component: AdminLogin // We'll create this later
    component: () => import('../views/admin/AdminLogin.vue'), // Lazy load login
    meta: { requiresGuest: true }, // Prevent logged-in admins from seeing login
  },
  {
    path: '/admin',
    // component: AdminLayout, // Wrapper for all admin pages
    component: () => import('../layouts/AdminLayout.vue'), // Lazy load layout
    meta: { requiresAuth: true }, // This route and children require admin auth
    children: [
      {
        path: '', // Default child for /admin
        name: 'admin-dashboard',
        component: () => import('../views/admin/AdminDashboard.vue') // Lazy load
      },
      {
        path: 'quizzes', // becomes /admin/quizzes
        name: 'admin-quiz-list',
        component: () => import('../views/admin/AdminQuizList.vue')
      },
      {
        path: 'quizzes/create', // Matches /admin/quizzes/create
        name: 'admin-quiz-create',
        component: () => import('@/views/admin/AdminQuizCreate.vue') // Lazy load
      },
      {
        path: 'quizzes/:id/edit',
        name: 'admin-quiz-edit',
        component: () => import('@/views/admin/AdminQuizEdit.vue'), // Lazy load
        props: true // Pass route params as props
      },
      {
        path: 'quizzes/:id/responses',
        name: 'admin-quiz-responses',
        component: () => import('@/views/admin/AdminResponseList.vue'), // Lazy load
        props: true // Pass route params as props
      },
      {
        path: 'quizzes/:id/credentials',
        name: 'admin-quiz-credentials',
        component: () => import('@/views/admin/AdminCredentialManager.vue'), // Lazy load
        props: true // Pass route params as props
      },
      {
        path: 'responses/:responseId',
        name: 'AdminResponseDetails',
        component: () => import('@/views/admin/AdminResponseDetails.vue'), // Lazy load the new component
        props: true // Pass responseId as prop
      }
      // Add other admin child routes here
    ],
  },
  // Responder routes
  {
    path: '/login',
    name: 'responder-login',
    component: () => import('../views/responder/ResponderLogin.vue'),
    meta: { requiresResponderGuest: true }, // Prevent logged-in responders from seeing login
  },
  {
    path: '/quiz/:id',
    name: 'quiz-taker',
    component: () => import('../views/responder/QuizTaker.vue'),
    props: true, // Pass route params as props
    meta: { requiresResponderAuth: true }, // This route requires responder auth
  },
  {
    path: '/quiz/:id/result',
    name: 'quiz-result',
    component: () => import('../views/responder/QuizResult.vue'),
    props: true, // Pass route params as props
    meta: { requiresResponderAuth: true }, // This route requires responder auth
  },

  // Catch-all 404
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: { template: '<div>404 Not Found</div>' },
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

// --- Navigation Guards ---
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore();
  
  // Log navigation for debugging
  console.log(`Navigating from ${from.fullPath} to ${to.fullPath}`);
  
  // Admin auth check
  if (to.meta.requiresAuth && !authStore.isAdminAuthenticated) {
    console.log('Admin authentication required, redirecting to login');
    next({ name: 'admin-login' });
    return;
  }
  
  // Admin guest check (prevent logged-in admins from seeing login)
  if (to.meta.requiresGuest && authStore.isAdminAuthenticated) {
    console.log('Already authenticated as admin, redirecting to dashboard');
    next({ name: 'admin-dashboard' });
    return;
  }
  
  // Responder auth check
  if (to.meta.requiresResponderAuth && !authStore.isResponderAuthenticated) {
    console.log('Responder authentication required, redirecting to login');
    next({ name: 'responder-login' });
    return;
  }
  
  // Responder guest check (prevent logged-in responders from seeing login)
  if (to.meta.requiresResponderGuest && authStore.isResponderAuthenticated) {
    console.log('Already authenticated as responder, redirecting to quiz');
    const quizId = authStore.responderQuizId;
    if (quizId) {
      next({ name: 'quiz-taker', params: { id: quizId.toString() } });
    } else {
      // If for some reason we don't have the quiz ID, just go to home
      next({ name: 'home' });
    }
    return;
  }
  
  // Allow navigation
  next();
});

export default router;
