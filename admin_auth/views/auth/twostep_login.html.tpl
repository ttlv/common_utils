<style>
  .title {
    color: white;
    width: 100%;
    text-align: center;
  }
  .login-page {
    width: 360px;
    padding: 8% 0 0;
    margin: auto;
  }
  .phone {
   float: left;
   margin-bottom: 12px;
   margin-left: 1px;
   color: #444040;
  }
  .form {
    position: relative;
    z-index: 1;
    background: #FFFFFF;
    max-width: 360px;
    margin: 0 auto 100px;
    padding: 45px;
    text-align: center;
    box-shadow: 0 0 20px 0 rgba(0, 0, 0, 0.2), 0 5px 5px 0 rgba(0, 0, 0, 0.24);
  }
  .form input {
    font-family: "Roboto", sans-serif;
    outline: 0;
    background: #f2f2f2;
    width: 100%;
    border: 0;
    margin: 0 0 15px;
    padding: 15px;
    box-sizing: border-box;
    font-size: 14px;
  }
  .form button {
    font-family: "Roboto", sans-serif;
    text-transform: uppercase;
    outline: 0;
    background: #0277bd;
    width: 100%;
    border: 0;
    padding: 15px;
    color: #FFFFFF;
    font-size: 14px;
    -webkit-transition: all 0.3 ease;
    transition: all 0.3 ease;
    cursor: pointer;
  }
  .form button:hover,.form button:active,.form button:focus {
    background: #0277b0;
  }
  .form .message {
    margin: 15px 0 0;
    color: #b3b3b3;
    font-size: 12px;
  }
  .form .message a {
    color: #4CAF50;
    text-decoration: none;
  }
  .form .register-form {
    display: none;
  }
  .container {
    position: relative;
    z-index: 1;
    max-width: 300px;
    margin: 0 auto;
  }
  .container:before, .container:after {
    content: "";
    display: block;
    clear: both;
  }
  .container .info {
    margin: 50px auto;
    text-align: center;
  }
  .container .info h1 {
    margin: 0 0 15px;
    padding: 0;
    font-size: 36px;
    font-weight: 300;
    color: #1a1a1a;
  }
  .container .info span {
    color: #4d4d4d;
    font-size: 12px;
  }
  .container .info span a {
    color: #000000;
    text-decoration: none;
  }
  .container .info span .fa {
    color: #EF3B3A;
  }
  .error {
    padding-bottom: 15px;
    color: red;
  }
  body {
    background: #03a9f4; /* fallback for old browsers */
    font-family: sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }
</style>
<div class="login-page">
  <h1 class="title">Admin Login</h1>
  <div class="form">
    {{if .error}}
      <div class="error">{{.error}}</div>
    {{end}}
    <form class="login-form" action="{{mountpathed "login"}}" method="POST">
      <div class="phone">验证码已发送至: {{.phone}}</div>
      <input type="code" class="form-control" name="code" placeholder="请输入收到的手机验证码"><br />
      <input type="hidden" name="{{.xsrfName}}" value="{{.xsrfToken}}" />
      <input type="hidden" name="{{.primaryID}}" value="{{.primaryIDValue}}" />
      <input type="hidden" name="token" value="{{.token}}" />
      <button>确定</button>
    </form>
  </div>
</div>
