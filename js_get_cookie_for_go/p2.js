var fs = require('fs');
var system = require('system');
var args = system.args;
var uname1, pwd1, uname2, pwd2;
if (args.length < 5) {
    console.log('Missing parameter!');
    phantom.exit();
} else {
    uname1 = args[1];
    pwd1   = args[2];
    uname2 = args[3];
    pwd2   = args[4];
}
var cookies_path = './cookies/'+uname1+'-'+uname2+'.json';
var page = require('webpage').create();
var cookies = '';
page.settings.userAgent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.87 Safari/537.36';
page.settings.resourceTimeout = 15000;
page.loadImages = false;


page.open('https://glogin.rms.rakuten.co.jp/', function (status) {
    var contentHtml_1 = page.content;
    if(status == 'success' && contentHtml_1 && contentHtml_1.indexOf('R-Loginにログインする') > 0){
        page.evaluate(function(uname, pwd) {
            document.querySelector('#rlogin-username-ja').value = uname;
            document.querySelector('#rlogin-password-ja').value = pwd;
            document.querySelector('.rf-button-primary').click();
            console.log('step1ok');

        }, uname1, pwd1);
    }else{
        exit_program()
    }
});

var times = 0;
page.onConsoleMessage = function(msg, lineNum, sourceId) {
    if (msg === 'step1ok') {
        var timer1 = setInterval(function () {
            var pageHtml_2 = page.content;
            if(pageHtml_2 && pageHtml_2.indexOf('で認証済み') > 0){
                times = 0;
                clearInterval(timer1);

                page.evaluate(function(uname, pwd) {
                    document.querySelector('#rlogin-username-2-ja').value = uname;
                    document.querySelector('#rlogin-password-2-ja').value = pwd;
                    document.querySelector('.rf-button-primary').click();
                    console.log('step2ok');
                }, uname2, pwd2);
            }else{
                times = times + 1;
                if(times > 3){
                    times = 0;
                    clearInterval(timer1);
                    exit_program('First login failed');
                }
            }
        }, 5000);
    } else if (msg === 'step2ok') {
        var timer2 = setInterval(function () {
            var pageHtml_3 = page.content;
            if(pageHtml_3 && pageHtml_3.indexOf('お気をつけください') > 0){
                times = 0;
                clearInterval(timer2);

                page.evaluate(function() {
                    document.querySelector('.rf-button-primary').click();
                    console.log('step3ok');
                });
            }else{
                times = times + 1;
                if(times > 3){
                    times = 0
                    clearInterval(timer2);
                    exit_program('Second login failed');
                }
            }
        }, 5000);
    } else if (msg === 'step3ok') {
        var timer3 = setInterval(function () {
            var pageHtml_4 = page.content;
            if(pageHtml_4 && pageHtml_4.indexOf('楽天からの重要なお知らせ') > 0){
                times = 0;
                clearInterval(timer3);

                page.evaluate(function() {
                    document.querySelector('#com_gnavi0303').click();
                    console.log('step4ok');
                });
            }else{
                times = times + 1;
                if(times > 3){
                    times = 0
                    clearInterval(timer3);
                    exit_program('Login Confirm failed');
                }
            }
        }, 5000);
    } else if(msg === 'step4ok'){
        var timer4 = setInterval(function () {
            var pageHtml_5 = page.content;
            if(pageHtml_5 && pageHtml_5.indexOf('R-Datatool アクセス分析') > 0){
                times = 0;
                clearInterval(timer4);
                //var report_url = 'https://rdatatool.rms.rakuten.co.jp/access/?menu=pc&evt=RT_P03_01&stat=1&owin=';
                //retry2open(report_url, 3)
                page.evaluate(function() {
                    document.querySelector('#mm_sub0303_08').click();
                    console.log('step5ok');
                });
            }else{
                times = times + 1;
                if(times > 3){
                    times = 0
                    clearInterval(timer4);
                    exit_program('Opening R-Datatool page failed');
                }
            }
        }, 5000);
    }else if(msg === 'step5ok'){
        var timer5 = setInterval(function () {
            var pageHtml_6 = page.content;
            if(pageHtml_6 && pageHtml_6.indexOf('商品ページランキング') > 0){
                times = 0;
                clearInterval(timer5);

                console.log('get_cookies_success')
                cookies = JSON.stringify(page.cookies);
                exit_program(cookies);
            }else{
                times = times + 1;
                if(times > 3){
                    times = 0
                    clearInterval(timer5);
                    exit_program('Opening R-Datatool report page failed');
                }
            }
        }, 5000)
    }
    console.log('CONSOLE: ' + msg + ' (from line #' + lineNum + ' in "' + sourceId + '")');
};

setTimeout(function() {
    exit_program('Timeout');
}, 80000);

function exit_program(msg) {
    if(msg) {
        console.log('Message: ', msg);
    }
    fs.write(cookies_path, cookies, 'w');
    page.close();
    phantom.exit();
    return
}

/*
function retry2open(page_url, retries) {
    if(retries > 0){
        page.open(page_url, function (status) {
            var pageHtml_6 = page.content;
            console.log(status, pageHtml_6.length)
            cookies = pageHtml_6;
            exit_program('====================');
            if(status == 'success' && pageHtml_6 && pageHtml_6.indexOf('商品ページランキング') > 0){
                //cookies = JSON.stringify(page.cookies);
                exit_program(cookies);
            }else{
                retries --
                retry2open(page_url, retries)
            }
        });
    }else{
        exit_program('opening report page failed');
    }
    return
}*/
