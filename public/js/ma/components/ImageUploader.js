import $ from 'jquery'

export default class ImageUploader {
    constructor(api, selector) {
        this._api = api
        this._$file = $(selector)
    }

    upload(path){
        var formData = new FormData()
        $.each(this._$file[0].files, function(i, file){
            formData.append('file', file)
        })
        formData.append('path', path)

        return $.ajax({
            url: this._api,
            method: 'post',
            dataType: 'json',
            data: formData,
            processData: false,
            contentType: false
        })
    }

    validate(){
        $.each(this._$file[0].files, function(i, file){
            // file size valdation by 2MB
            if(file.size > 2000000){
                alert('2MB以上のファイルはアップロードできません')
                return false            
            }
        })
        return true
    }

    getFileName(){
      return this._$file.val().replace(/\\/g, '/').replace(/.*\//, '')
    }
}