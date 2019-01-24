Date.prototype.getLocalTime = function() {
    h = (h = this.getHours()) < 10 ? ('0' + h) : h;
    m = (m = this.getMinutes()) < 10 ? ('0' + m) : m;
    s = (s = this.getSeconds()) < 10 ? ('0' + s) : s;
    return h + ":" + m + ":" + s;
}

$.routes.add('/{id:int}/', 'issueNumRoute', function() {
    $('#picker_select').val(this.id)
    $('#picker_select').trigger("change");
});

$(document).ready(function(){

    $('#totop').click(function(e){
        e.preventDefault();
        $("#messages").prop({ scrollTop: 0 });
    });

    $('#tobottom').click(function(e){
        e.preventDefault();
        $("#messages").prop({ scrollTop: $("#messages")[0].scrollHeight });
    });

    $.jPlayer.timeFormat.showHour = true;

    $("#jquery_jplayer_1").jPlayer({
        solution: "html, flash",
        swfPath: "//cdnjs.cloudflare.com/ajax/libs/jplayer/2.9.2/jplayer/jquery.jplayer.swf",
        supplied: "mp3",
        preload: "metadata",
        wmode: "window",
        ready: function () {
            $(this).jPlayer("setMedia", {
                mp3: "http://91.213.196.99:8181/stream"
            });
        },
        timeupdate: function (event) {
            var last_offset = parseInt($('#last_offset').text());
            var cur_offset = event.jPlayer.status.currentTime;
            if (last_offset == cur_offset) { 
                //$('#last_offset').text(cur_offset.toString());
                return;
            }
            var start_time = parseInt($("#start_time").text());
            var elem = $('#messages');
            if (cur_offset > last_offset) {
              if (elem[0].scrollHeight - elem.scrollTop() == elem.outerHeight()) {
                for(var i=last_offset; i <= cur_offset; i++) {
                  $(".t"+(start_time+i).toString()).css("display","block");
                }
                $("#messages").prop({ scrollTop: $("#messages").prop("scrollHeight") });
              } else {
                for(var i=last_offset; i <= cur_offset; i++) {
                  $(".t"+(start_time+i).toString()).css("display","block");
                }
              }
            } else {
              for(var i=last_offset; i >= cur_offset; i--) {
                $(".t"+(start_time+i).toString()).css("display","none");
              }
            }
            $('#last_offset').text(cur_offset.toString());
        },
        ended: function () {
            $(".message").css("display","block");
        }
    });

    $('#picker_select').change(function(){
        var num = $(this).val();

        $('.message').remove();
        $('#live_chat').remove();
        //$('#numbers').css("display","inline-table");
        $('#loaded').text('0');
        $('#total').text('0');
        $('#fork_me').css("display","block");

        if (num == "--") {
            return;
        }

        location.href = "/#/"+num

        if (num == "live") {
            //$("#jquery_jplayer_1").jPlayer("setMedia",{mp3: "http://91.213.196.99:8181/stream"});
            $("#jquery_jplayer_1").jPlayer("setMedia",{mp3: "http://pf.volna.top/PilotBy48"});
            //$('#messages').append('<iframe src="http://chat.radio-t.com/" id="live_chat"><p>Your browser does not support iframes.</p></iframe>');
            $('#messages').append('<div class="message" >Try to implement Giiter integration later.</div>');
            $('.message').css("display","block");
            $('#fork_me').show();

            return;
        }

        $.getJSON("data/"+num+"/desc.json", function (data) {
            console.log(data)
            $('#start_time').text(data.start_time);
            $("#jquery_jplayer_1").jPlayer("setMedia",{mp3: data.url});
            $('#total').text(data.chat_n);
        });

        $.getJSON("data/"+num+"/chat.json", function (data) {
            var items = [];
            $.each( data.chat, function( key, message ) {
                $('#messages').append('<div class="message t' + message.datetime + '"><p><span class="date">' + (new Date(message.datetime*1000).getLocalTime()) + '</span><span class="nickname ' + message.author_type + '">' + message.author_name + '</span><span class="text">' + message.text + '</span></p></div>');
                $('#loaded').text(parseInt($('#loaded').text())+1);
            });

            var start_time = parseInt($('#start_time').text());
            for (var i = start_time - 10800; i < start_time; i++ ) {
                $('.t' + i).css("display","block");
            }
            $('#tobottom').click();
        });
    });

    $('#picker_select').trigger("change");
});