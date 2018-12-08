import Vue from 'vue';
import VueRouter from 'vue-router';
import VueI18n from 'vue-i18n';
// import socketio from 'socket.io';
import VueSocketIO from 'vue-socket.io';
import io from 'socket.io-client'
// require('default-passive-events')

import App from './App';
import Index from './components/public/Index';
import Movie from './components/private/movie/Movie';
import MovieGallery from './components/private/MovieGallery';
import ShowCategory from './components/private/ShowCategory';
import Search from './components/private/search/Search';
import NotFound from './components/NotFound';
import messages from './static/data/locales';
import RecoverPasswordHash from './components/public/RecoverPasswordHash';
import RecoverPassword from './components/public/RecoverPassword';
import RegisterForm from './components/public/RegisterForm';
import LoginForm from './components/public/LoginForm';
import { getCurLanguage } from './static/js/language';
import auth from './static/js/auth';

export const socket = io(process.env.API_DOMAIN_NAME, {
  transports: ['websocket']
});

const queryString = require('query-string');

Vue.config.productionTip = false;

Vue.filter('uppercase', value => (!value ? '' : value.toString()
  .toUpperCase()));

Vue.use(VueRouter);
Vue.use(VueI18n);

socket.on('reconnect', e => {
  // console.log("reconnect socket")
})
socket.on('disconnect', e => {
  // console.log("disconnect socket")
})
// Vue.use(VueSocketIO, SocketInstance);

// SocketInstance.on('connect', (e) => null);
// SocketInstance.on('disconnect', (e) => ull);


const EventBus = new Vue({});

Object.defineProperties(Vue.prototype, {
  $_bus: {
    get() {
      return EventBus;
    },
  },
  $_socket: { get: function () { return socket } },
});


const router = new VueRouter({
  mode: 'history',
  routes: [
    {
      path: '/',
      component: Index,
      children: [
        {
          path: '',
          component: LoginForm,
          meta: {
            title: 'Home',
            requiresAuth: false,
          },
          name: 'index',
        },
        {
          path: '/recover',
          component: RecoverPassword,
          meta: {
            title: 'Recover Password',
            requiresAuth: false,
          },
          name: 'recover',
        },
        {
          path: '/recover/:hash',
          component: RecoverPasswordHash,
          meta: {
            title: 'Recover Password',
            requiresAuth: false,
          },
          name: 'recover_hash',
        },
        {
          path: '/register',
          component: RegisterForm,
          meta: {
            title: 'Register',
            requiresAuth: false,
          },
          name: 'register',
        },
      ],
    },
    {
      path: '/gallery',
      component: MovieGallery,
      meta: {
        title: 'Movie Gallery',
        requiresAuth: true,
      },
    },
    {
      path: '/movie/:id(tt\\d+|\\d+)', // ttdddd... -> imdb / ddd... -> tmdb
      component: Movie,
      meta: {
        title: 'Movie',
        requiresAuth: true,
      },
    },
    {
      path: '/category/:name([a-z-]+)',
      component: ShowCategory,
      meta: {
        title: 'Category',
        requiresAuth: true,
      },
    },
    {
      path: '/search',
      component: Search,
      meta: {
        title: 'Search',
        requiresAuth: true,
      },
    },
    {
      path: '/404',
      component: NotFound,
      meta: {
        title: 'Page not found',
        requiresAuth: false,
      },
    },
  ],
});

const i18n = new VueI18n({
  locale: getCurLanguage(),
  messages,
});

router.beforeEach((to, from, next) => {
  document.title = `${to.meta.title} - Hypertube`;
  if (!to.matched.length) {
    return next('/404');
  }
  // Check if the route need an authentication
  if (to.matched.some(record => record.meta.requiresAuth)) {
    const parsed = queryString.parse(location.search);
    if (parsed.token) {
      // If the url contain a token store it in localstorage and update url
      window.history.pushState({}, '', to.path);
      auth.storeToken(parsed.token);
      return next(({
        path: to.path,
      }));
    } else if (!auth.isLoggedIn()) {
      return next({
        path: '/',
        query: { redirect: to.fullPath },
      });
    }
  } else if (auth.isLoggedIn()) {
    return next({
      path: '/gallery',
    });
  }
  return next();
});

new Vue({
  router,
  i18n,
  render: h => h(App),
}).$mount('#app');
