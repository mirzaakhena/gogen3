import {createRouter, createWebHistory} from 'vue-router'

const routes = [

    {
        path: '/',
        component: () => import('./PageWithSidebar.vue'),
        children: [
            {
                path: '/payment',
                redirect: '/payment/akhena',
                component: () => import('./yourpage/ViewTab.vue'),
                children: [
                    {
                        path: '/payment/mirza',
                        component: () => import('../usecase/getallyourusecasename/PageTable.vue'),
                    },
                    {
                        path: '/payment/akhena',
                        component: () => import('../usecase/runyourusecasename/PageButton.vue'),
                    },
                ],
            },
            {
                path: '/order',
                component: () => import('../usecase/getyourusecasename/PageButton.vue'),
            },
        ],
    },

]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router