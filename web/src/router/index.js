import Vue from "vue";
import Router from "vue-router";
import Layout from "@/layout";

Vue.use(Router);

export const constantRoutes = [
  {
    path: "/login",
    component: () => import("@/views/login/index"),
    hidden: true,
  },

  {
    path: "/404",
    component: () => import("@/views/404"),
    hidden: true,
  },

  {
    path: "/",
    component: Layout,
    redirect: "/developer",
    children: [
      {
        path: "developer",
        name: "Developer",
        component: () => import("@/views/developer/index"),
        meta: {title: "开发者账号管理", icon: "developer"},
      },
    ],
  },

  {
    path: "/ipa",
    component: Layout,
    children: [
      {
        path: "list",
        name: "IPA",
        component: () => import("@/views/ipa/index"),
        meta: {title: "应用管理", icon: "ipa"},
      },
    ],
  },

  {
    path: "/setting",
    component: Layout,
    children: [
      {
        path: "user",
        name: "Setting",
        component: () => import("@/views/setting/index"),
        meta: {title: "系统设置", icon: "setting"},
      },
    ],
  },

  {
    path: "github",
    component: Layout,
    children: [
      {
        path: "https://github.com/togettoyou/supersign",
        meta: {title: "项目地址", icon: "github"},
      },
    ],
  },

  // 404 page must be placed at the end !!!
  {path: "*", redirect: "/404", hidden: true},
];

const createRouter = () =>
  new Router({
    base: '/admin',
    // mode: 'history', // require service support
    scrollBehavior: () => ({y: 0}),
    routes: constantRoutes,
  });

const router = createRouter();

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter();
  router.matcher = newRouter.matcher; // reset router
}

export default router;
