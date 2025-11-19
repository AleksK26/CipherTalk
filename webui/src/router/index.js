import { createRouter, createWebHashHistory } from "vue-router";
import LoginView from "../views/LoginView.vue";
import HomeView from "../views/HomeView.vue";
import SearchPeopleView from "../views/SearchPeopleView.vue";
import ChatView from "../views/ChatView.vue";
import ProfileView from "../views/ProfileView.vue";
import GroupsView from "../views/GroupsView.vue"
import GroupCreateView from "../views/GroupCreateView.vue"
import GroupEditView from "../views/GroupEditView.vue"
import NotFoundView from "../views/NotFoundView.vue";
import GroupView from "../views/GroupView.vue";

const routes = [
  { path: "/", component: LoginView },
  { path: "/home", component: HomeView },
  { path: "/search", component: SearchPeopleView },
  { path: "/conversations/:uuid", name: "ChatView", component: ChatView, props: true },
  { path: "/me", component: ProfileView},
  { path: "/groups", component: GroupsView},
  { path: "/new-group", component: GroupCreateView},
  { path: "/groups/:uuid", name: "GroupEditView", component: GroupEditView, props: true},
  { path: "/:pathMatch(.*)*", name: "NotFound", component: NotFoundView },
  { path: "/group/:groupId", name: "Group", component: GroupView, meta: { requiredAuth: true } }
];

const router = createRouter({
  history: createWebHashHistory(),
  routes,
});

const publicPages = ["/"];
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("token");
  if (!publicPages.includes(to.path) && !token) {
    next("/");
  } else {
    next();
  }
});

export default router;