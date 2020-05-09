$(document).ready(function(){
    $("body").on("click", ".js-export", function(i, e) {
      if (!$(this).hasClass("disabled")) {
        var timestamp = Date.now();
        var _this = $(this);
        $(this).addClass("disabled");
        var $loading = $('<div class="js-loading" style="text-align: center;"><div class="mdl-spinner mdl-js-spinner is-active qor-layout__bottomsheet-spinner"></div></div>');
        $loading.appendTo($('.readonly-filters')).trigger('enable.qor.material');
        $.fileDownload($(this).attr("href") + "?timestamp=" + timestamp + "&" + $(".readonly-filters form").serialize(), {
            cookieValue: String(timestamp),
            successCallback: function(url) {
              _this.removeClass("disabled");
              $(".js-loading").remove();
            },
            failCallback: function (responseHtml, url, error) {
              alert($(responseHtml).text());
              _this.removeClass("disabled");
              $(".js-loading").remove();
            },
            abortCallback: function (url) {
              _this.removeClass("disabled");
              $(".js-loading").remove();
            },
        });
      }
      return false;
    });
});
