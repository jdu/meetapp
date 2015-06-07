import $ from 'jquery'
import util from '../util'

export default class Discussion {
    constructor() {
        this._$wrap = $('.ma-discussion-wrap')
        this._$textarea = $('#ma_detail_discussion_comment')
        $('.ma_detail_discussion_comment_btn').on('click', () => {
            if(util.getUserInfo()){
                this.post()
            }else{
                // move to login for anonymous
                location.href = "/login"
            }
        })
        var that = this
        $('.ma-profile-comment-delete-btn').on('click', function(){
            if(window.confirm('このコメントを削除してもよろしいでしょうか？')){
                var id = $(this).data('discussion-id')
                if(id){
                    that.deleteDiscussion(id)
                }
            }
        })
    }

    deleteDiscussion(id) {
        $.ajax({
            url: '/u/api/app/discussion/' + id,
            type: 'delete'
        }).done((res) => {
            location.reload()
        }).fail(() => {
            alert("Error")
        })
    }

    post() {
        var params = this.getParams()
        //  validation
        var result = this.validate(params)
        if(result.error){
            alert(result.message)
            return
        }
        $.ajax({
            url: '/u/api/app/discussion',
            type: 'post',
            contentType:"application/json; charset=utf-8",
            dataType: 'json',
            data: JSON.stringify(params)
        }).done((res) => {
            location.reload()
        }).fail(() => {
            alert("Error")
        })
    }

    validate(params) {
        if($.trim(params.discussionInfo.message) === ""){
            return {"error": true, "message": "コメントが入力されていません"}
        }
        return {"error": false}
    }

    getParams() {
        return {
            appId: util.getAppDetailId(),
            discussionInfo: {
                userId: util.getUserInfo().ID,
                message: this._$textarea.val()
            }
        }
    }
}
