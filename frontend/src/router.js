import { createRouter, createWebHistory } from 'vue-router';
import AuthForm from './components/AuthForm.vue';
import Home from './views/Home.vue';

const routes = [
  {
    path: '/auth',
    name: 'Auth',
    component: AuthForm,
  },
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  // Add more routes here as needed
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
