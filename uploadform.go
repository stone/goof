package main


const uploadform = `
<!DOCTYPE html>

<html>
<head>
  <meta name="generator" content=
  "HTML Tidy for Linux (vers 25 March 2009), see www.w3.org">

  <title>goof - upload</title>
  <style type="text/css">
  body{
	font-family:"Lucida Grande", "Lucida Sans Unicode", Verdana, Arial, Helvetica, sans-serif;
	font-size:12px;
  }
  p, h1, form, button{border:0; margin:0; padding:0;}
	.spacer{clear:both; height:1px;}
	.myform{
	margin:0 auto;
	width:400px;
	padding:14px;
  }
	#stylized{
	border:solid 2px #aaa;
	background:#e9e9e9;
  }
	#stylized h1 {
	font-size:14px;
	font-weight:bold;
	margin-bottom:8px;
  }
  #stylized p{
	font-size:11px;
	color:#666666;
	margin-bottom:20px;
	border-bottom:solid 1px #b7ddf2;
	padding-bottom:10px;
  }
  #stylized label{
	display:block;
	font-weight:bold;
	text-align:right;
	width:140px;
	float:left;
  }
  #stylized .small{
	color:#666666;
	display:block;
	font-size:11px;
	font-weight:normal;
	text-align:right;
	width:140px;
  }
  #stylized input{
	float:left;
	font-size:12px;
	padding:4px 2px;
	border:solid 1px #aacfe4;
	width:200px;
	margin:2px 0 20px 10px;
  }
  #stylized button{
	clear:both;
	margin-left:60px;
	width:300px;
	height:31px;
  }
  .button {
	display: inline-block;
	zoom: 1; /* zoom and *display = ie7 hack for display:inline-block */
	*display: inline;
	vertical-align: baseline;
	margin: 0 2px;
	outline: none;
	cursor: pointer;
	text-align: center;
	text-decoration: none;
	font: 14px/100% Arial, Helvetica, sans-serif;
	padding: .5em 2em .55em;
	text-shadow: 0 1px 1px rgba(0,0,0,.3);
	-webkit-border-radius: .5em;
	-moz-border-radius: .5em;
	border-radius: .5em;
	-webkit-box-shadow: 0 1px 2px rgba(0,0,0,.2);
	-moz-box-shadow: 0 1px 2px rgba(0,0,0,.2);
	box-shadow: 0 1px 2px rgba(0,0,0,.2);
  }
  .button:hover {
	text-decoration: none;
  }
  .button:active {
	position: relative;
	top: 1px;
  }
  .bigrounded {
	-webkit-border-radius: 2em;
	-moz-border-radius: 2em;
	border-radius: 2em;
  }
  .medium {
	font-size: 12px;
	padding: .4em 1.5em .42em;
  }
  .small {
	font-size: 11px;
	padding: .2em 1em .275em;
  }
  .gray {
	color: #e9e9e9;
	border: solid 1px #555;
	background: #6e6e6e;
	background: -webkit-gradient(linear, left top, left bottom, from(#888), to(#575757));
	background: -moz-linear-gradient(top,  #888,  #575757);
	filter:  progid:DXImageTransform.Microsoft.gradient(startColorstr='#888888', endColorstr='#575757');
  }
  .gray:hover {
	background: #616161;
	background: -webkit-gradient(linear, left top, left bottom, from(#757575), to(#4b4b4b));
	background: -moz-linear-gradient(top,  #757575,  #4b4b4b);
	filter:  progid:DXImageTransform.Microsoft.gradient(startColorstr='#757575', endColorstr='#4b4b4b');
  }
  .gray:active {
	color: #afafaf;
	background: -webkit-gradient(linear, left top, left bottom, from(#575757), to(#888));
	background: -moz-linear-gradient(top,  #575757,  #888);
	filter:  progid:DXImageTransform.Microsoft.gradient(startColorstr='#575757', endColorstr='#888888');
  }
  </style>
</head>

<body>
  <div id="stylized" class="myform">
    <form action="/upload" method="post" id="uploadform" enctype="multipart/form-data">
      <h1>Upload</h1>

      <p>Powered by goof</p>
		<label>File: <span class="small">Choose file to upload</span></label>
		<input type="file" id="fileinput" multiple="true" name="file">
		<button class="button gray" type="submit">Send in the binary chaos</button>
		<div class="spacer"></div>
	</form>
	
	<a target="_blank" href="https://github.com/stone/goof">goof</a>
  </div>
</body>
</html>
`
