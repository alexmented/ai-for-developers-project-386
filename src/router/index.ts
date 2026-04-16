import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'landing',
    component: () => import('@/pages/LandingPage.vue'),
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('@/pages/AdminPage.vue'),
  },
  {
    path: '/:ownerSlug',
    name: 'owner-event-types',
    component: () => import('@/pages/PublicEventTypesPage.vue'),
  },
  {
    path: '/:ownerSlug/:eventTypeId',
    name: 'owner-booking',
    component: () => import('@/pages/PublicBookingPage.vue'),
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
