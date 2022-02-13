import Vue from 'vue'
import VueRouter from 'vue-router'
import Login from '@/views/Login'
import Layout from '@/layout'

Vue.use(VueRouter)

const routes = [
    {
        path: '/login',
        name: 'Login',
        component: Login,
        meta: {
            allowAnonymous: true
        }
    },
    {
        path: '/',
        name: 'layout',
        redirect: '/overview',
        component: Layout,
        children: [
            {
                path: '/overview',
                name: 'Overview',
                component: () => import(/* webpackChunkName: "overview" */ '../views/Overview.vue')
            },
            {
                path: '/domains',
                name: 'Domains',
                component: () => import(/* webpackChunkName: "domains" */ '../views/Domains.vue')
            },
            {
                path: '/records',
                name: 'Records',
                component: () => import(/* webpackChunkName: "records" */ '../views/Records.vue')
            }
        ]
    }
]

const router = new VueRouter({
    mode: 'history',
    base: process.env.BASE_URL,
    routes
})

export default router
