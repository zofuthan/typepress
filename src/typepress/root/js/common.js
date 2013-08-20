(function(jQuery, window) {
    var bp = window.App = window.App || {};
    var i, j, k, s, o, a;
    var autoid = 0;
    init()
    if (!jQuery) {
        return alert("Where is jQuery/Zepto ?")
    }

    jQuery("form").submit(onSubmit);
    jQuery("label").each(function() {
        if (jQuery(this).attr('for') != undefined) {
            return
        }
        i = jQuery(this).next("input").eq(0);
        if (i.size() == 1) {
            if (i.attr("id") == "") {
                autoid++;
                i.attr('id', aid());
                jQuery(this).attr('for', aid());
                return
            }
            jQuery(this).attr('for', i.attr("id"));
        }
    });
    var act = {};

    function aid() {
        return 'aid_' + autoid;
    }

    function alert(msg) {
        window.alert(msg)
        return false
    }

    function notify(msg) {
        window.alert(msg)
        return false
    }

    function init() {
        jQuery.extend(bp, {
            act: act,
            alert: alert,
            notify: notify,
            onSubmit: onSubmit,
            aQ: aQ,
            ajaxDone: ajaxDone,
            ajaxFail: ajaxFail,
            anv2o: anv2o,
            jsonParse: jsonParse,
            pagIn: pagIn
        })
    }

    function onSubmit() {
        var data;
        var url = jQuery(this).attr('action')||"";
        if (url!="" && url.slice(-1) == "#") {
            return true
        }
        var type = this.method.toUpperCase();
        var dst = jQuery(this).closest('[id]');
        // act 检查
        var call = act[url + ":onsubmit"]
        if (call && !call(type, this)) {
            return false
        }
        if (jQuery(this).find('[type=file]').size()) {
            data = new FormData(this);
        } else {
            data = jQuery(this).serializeArray()
            data = anv2o(data)
            if (data.User_login != undefined) {
                data.User_login = md5(data.User_login)
            }
            if (data.user_login != undefined) {
                data.user_login = md5(data.user_login)
            }
            if (data.user_pass != undefined) {
                data.user_pass = md5(data.user_pass)
            }
            if (data.User_pass != undefined) {
                if (data.confirm != undefined && data.confirm != data.User_pass) {
                    return notify.call(jQuery(this, "[name=User_pass]"), "两次密码录入不一样")
                }
                delete data.confirm
                data.User_pass = md5(data.User_pass)
            }
        }
        return aQ({
            url: url,
            type: type
        }, {
            dst: dst
        }, data);

    }

    //ajax queue

    function aQ(url, extData, data) {
        var a = (typeof url == 'string') ? {
            url: url
        } : url;
        if (!a.data) {
            a.data = data
        }
        if (!a.extData) {
            a.extData = extData
        }
        if (a.data instanceof FormData) {
            a.processData = false;
            a.contentType = false;
        }
        a.success = a.success || ajaxDone;
        a.error = a.error || ajaxFail;
        jQuery.ajax(a);
        return false;
    }

    function ajaxDone(txt, stat, xhr) {
        var ext = this.extData || {};
        var t = ext && ext.apply || null;
        txt = txt.trimLeft();
        var method = this.type;
        // 简化智能处理
        switch (txt[0]) {
            case ';': //javascript
                (new Function(txt)).apply(t, ext.args || []);
                return;
            case '/': //url
                if (window.location.pathname == txt) {
                    window.location.reload();
                } else {
                    window.location.href = txt;
                }
                return;
            case '<': //html
                dst = jQuery(ext.dst || '#con')
                if (!dst.size())
                    dst = jQuery('#con');
                var tag = "";
                for (var i = 1;; i++) {
                    if (txt[i] == " " || txt[i] == ">")
                        break;
                    tag += txt[i]
                }
                switch (tag) {
                    case 'tbody':
                        dst = dst.find('table:first').append(txt);
                        break;
                    case 'tr':
                        dst = dst.find('table:first>tbody:last').append(txt);
                        break;
                    case 'div':
                        dst.html(txt);
                        break;
                    default:
                        return notify(tag);
                }
                clearTextNode(dst.get(0));
                //ajax来的url需要继续让ajax支持
                pagIn(this.url.slice(1), dst);
                return;
            case '{': //json数据有可能带有模板源码
            case '[':
                var json = txt.indexOf("\n");
                var call = act[this.url];
                if (call) {
                    if (json == -1) {
                        json = txt;
                        txt = "";
                    } else {
                        json = txt.slice(0, json);
                        txt = txt.slice(json.length + 1)
                    }
                    try {
                        json = JSON.parse(json);
                    } catch (e) {
                        json = false;
                    }
                    if (json) {
                        call(this.type, json, txt);
                    } else {
                        notify("服务器返回了错误的JSON格式:" + this.url)
                    }
                    break;
                }
                notify("无法处理:" + this.url)
                return;
            default:
                var loc = xhr.getResponseHeader("Location")
                if (loc) {
                    if (window.location.pathname == loc) {
                        window.location.reload();
                    } else {
                        window.location.href = loc;
                    }
                    break;
                }
                notify(txt || "糟糕,开发者忘记实现这个功能了!");
        }
    }

    function ajaxFail(xhr, stat, info) {
        notify("发生错误:" + (xhr.responseText || info));
    }
    //[{name:"",value:""}] 结构转 {key:value}

    function anv2o(nv) {
        var ret = {};
        if (!jQuery.isArray(nv)) {
            nv = [nv];
        }
        for (var i = 0; i < nv.length; i++) {
            ret[nv[i].name] = nv[i].value;
        }
        return ret;
    }

    function jsonParse(json) {
        var o;
        try {
            o = JSON.parse(json);
        } catch (e) {
            return undefined;
        }
        return o;
    }
    //分页处理

    function pagIn(url, range) {
        range.find('.pagination').each(function() {
            var t = jQuery(this),
                l = t.data('paginlimit') || 0,
                b = jQuery('>b', this).eq(0),
                at = parseInt(b.text()) || 0,
                s = Math.max(1, at - 5),
                e = s + 9;
            if (!l || !at) return;
            url = url.replace(/([?&])s=\d+/, '$1').replace('&&', '&').replace('?&', '?').replace(/\?$/, '');
            var prefix;
            if (url.indexOf('?') == -1) {
                prefix = '<a href="' + url + '?s=';
            } else {
                prefix = '<a href="' + url + '&s=';
            }
            //绘制prev，，，
            if (s >= 2) {
                jQuery(b).before(prefix + Math.max(at - 5, 1) + '">上翻</a>');
            } else {
                jQuery(b).before('<b>上翻<b>');
            }
            //没有最后限制
            jQuery(b).after(prefix + (Math.max(at, 6) + 5) + '">下翻</a>');
            for (; s < at; s++) {
                jQuery(b).before(prefix + s + '">' + s + '</a>');
            }
            s = e;
            for (; s > at; s--) {
                jQuery(b).after(prefix + s + '">' + s + '</a>');
            }
            //setTimeout(function(){log(t.width());t.width(t.width()).addClass('mid');}, 0);

        });
    }
})(window.jQuery || window.Zepto, window)
