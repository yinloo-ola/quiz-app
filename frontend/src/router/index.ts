import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router';

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
      }
      // Add other admin child routes here
    ],
  },
  // Add other top-level routes here (e.g., responder quiz view)

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

// --- Navigation Guards (Placeholder) ---
// We will add guards later to check for authentication
router.beforeEach((to, from, next) => {
  // TODO: Implement auth checks based on meta fields (requiresAuth, requiresGuest)
  // For now, allow all navigation
  console.log(`Navigating from ${from.fullPath} to ${to.fullPath}`);
  next();
});

export default router;
