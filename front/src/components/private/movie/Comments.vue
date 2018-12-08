<template lang="html">
  <div class="comments-area">
    <user-profile :toggleUserProfil="toggleUserProfil" :userProfil="profilToShow" v-if="showUserProfil"></user-profile>
    <h2 class="box-title">{{ $t('commentLabel') }}</h2>
    <div v-if="!commentLoading">
      <div
        contenteditable="true"
        class="new-comment comment-box"
        @input="commentInput($event.target.innerText)"
        :placeholder="$t('commentInputPlaceholder')"
      >
      </div>
      <div class="comment-btn-box" v-if="canPushComment">
        <button
          class="comment-btn-box-cancel"
          @click="cancelCommentInput()"
        >
          {{ $t('cancelCommentLabel') }}
        </button>
        <button
          class="comment-btn-box-add"
          @click="addComment"
        >
          {{ $t('addCommentLabel') }}
        </button>
      </div>

      <div :key="index"
           class="comment-box"
           v-for="(comment, index) in comments"
           v-if="comments.length > 0"
      >
        <div class="comment-picture-frame">
          <img
            :src="comment.profile_picture"
            @click="toggleUserProfil(comment)"
            alt="comment user picture"
            class="adapt-picture-circle"
          >
        </div>
        <div class="comment-data">
          <p>
            <span class="comment-username" @click="toggleUserProfil(comment)">{{ comment.fullname }}</span>
            <span class="comment-ago">{{ comment.createdat }}</span>
          </p>
          <p class="comment">{{ comment.content }}</p>
        </div>
      </div>

      <div class="empty-comment-box" v-if="comments.length === 0">
        {{ $t('emptyFilmComments') }}
      </div>

    </div>
    <div v-else>
      <loading-animation></loading-animation>
    </div>
  </div>
</template>

<script>
import moment from 'moment';
import UserProfile from './UserProfile';
import LoadingAnimation from '../module/LoadingAnimation';
import api from '../../lib/api';
import auth from '../../../static/js/auth';
import swal from 'sweetalert';


export default {
  props: {
    filmId: String,
    comments: Array,
    commentLoading: Boolean
  },
  components: {
    UserProfile,
    LoadingAnimation,
  },
  data() {
    return {
      commentContent: '',
      showUserProfil: false,
      canPushComment: false,
      profilToShow: null
    };
  },
  methods: {
    timeSince(time) {
      
      return moment(time)
        .fromNow();
    },

    toggleUserProfil(comment) {
      if (this.showUserProfil) {
        this.showUserProfil = false;
        this.profilToShow = null;
      } else {
        this.showUserProfil = true;
        this.profilToShow = {
          picture: comment.profile_picture,
          local: comment.user_locale,
          fullname: comment.fullname
        }
      }
    },

    commentInput(content) {
      if (content === '') {
        this.commentContent = null;
        this.canPushComment = false;
      }
      else {
        this.commentContent = content;
        this.canPushComment = true;
      }
    },

    cancelCommentInput() {
      const editableDiv = document.querySelector('.new-comment')

      editableDiv.innerText = '';
      this.commentContent = '';
      this.canPushComment = false;
    },

    addComment() {
      if (this.commentContent) {
        api.addFilmComments(this.filmId, this.commentContent, auth.token())
        .then(res => {
          if (res.status === 200) {
            res.json()
              .then(data => {
                this.cancelCommentInput()
                this.$parent.updateComments(data)
              })
              .catch(error => swal({ text: this.$t('errorServer'), icon: 'error' }));
          } else {
            swal({ text: this.$t('errorServer'), icon: 'error' })
          }
        })
        .catch(err => swal({ text: this.$t('errorServer'), icon: 'error' }))
      }
    },
  },
};
</script>

<style scoped lang="sass">
  @import "Comments"
</style>
