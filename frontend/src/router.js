import { createRouter, createWebHistory } from 'vue-router';
import AuthForm from './components/AuthForm.vue';

const routes = [
  {
    path: '/auth',
    name: 'Auth',
    component: AuthForm,
  },
  // Add more routes here as needed
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
