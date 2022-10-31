import { createRouter, createWebHistory } from 'vue-router'
import MainPage from '@/components/MainPage.vue'
import KeyPage from "@/components/KeyPage";
import alreadyAuthenticated from "@/components/AlreadyAuthenticated";
import QuicktaskPage from "@/components/QuicktaskPage";

export default createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            component: MainPage,
        },
        {
            path: '/key',
            component: KeyPage,
        },
        {
            path: '/authenticated',
            name: 'authenticated',
            component: alreadyAuthenticated
        },
        {
            path: '/quicktask',
            name: 'quicktask',
            component: QuicktaskPage
        }
    ]
})