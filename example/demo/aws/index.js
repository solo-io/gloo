'use strict'

const aws = require('aws-sdk')
const s3 = new aws.S3()

exports.handler = (event, context, callback) => {
    // save data into s3
    if (event["email"]) {
        processForm(event, callback);
    } else {
        displayForm(callback);
    }
    
};

function displayForm(callback) {

    const html = `<!DOCTYPE html>
    <html>
      <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8"/>
        <meta charset="utf-8">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <link rel="shortcut icon" type="image/x-icon" href="/resources/images/favicon.png">
        <title>PetClinic :: a Spring Framework demonstration</title>
    
        <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/html5shiv/3.7.2/html5shiv.min.js"></script>
        <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
        <![endif]-->
    
        <link rel="stylesheet" href="/resources/css/petclinic.css"/>
      </head>
      <body>
      <nav class="navbar navbar-default" role="navigation">
          <div class="container">
              <div class="navbar-header">
                  <a class="navbar-brand" href="/"><span></span></a>
                  <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#main-navbar">
                      <span class="sr-only"><os-p>Toggle navigation</os-p></span>
                      <span class="icon-bar"></span>
                      <span class="icon-bar"></span>
                      <span class="icon-bar"></span>
                  </button>
              </div>
              <div class="navbar-collapse collapse" id="main-navbar">
                  <ul class="nav navbar-nav navbar-right">
      
                      <li>
                          <a href="">
                            <span class="glyphicon  glyphicon-null" aria-hidden="true"></span>
                            <span></span>
                          </a>
                      </li>
      
                      <li>
                          <a href="/" title="home page">
                            <span class="glyphicon  glyphicon-home" aria-hidden="true"></span>
                            <span>Home</span>
                          </a>
                      </li>
      
                      <li>
                          <a href="/owners/find" title="find owners">
                            <span class="glyphicon  glyphicon-search" aria-hidden="true"></span>
                            <span>Find owners</span>
                          </a>
                      </li>
      
                      <li>
                          <a href="/vets.html" title="veterinarians">
                            <span class="glyphicon  glyphicon-th-list" aria-hidden="true"></span>
                            <span>Veterinarians</span>
                          </a>
                      </li>
    
                      <li  class="active">
                        <a href="/contact.html" title="contact">
                          <span class="glyphicon  glyphicon-envelope" aria-hidden="true"></span>
                          <span>Contact</span>
                        </a>
                    </li>
                      <li>
                          <a href="/oups" title="trigger a RuntimeException to see how it is handled">
                            <span class="glyphicon  glyphicon-warning-sign" aria-hidden="true"></span>
                            <span>Error</span>
                          </a>
                      </li>
      
                  </ul>
              </div>
          </div>
      </nav>
      <div class="container-fluid">
        <div class="container xd-container">
        <h2>Contact</h2>
        <div class="container">
        <div id="notification" class="row"></div>
        <form method="post" action="/contact" id="contact">
        <div class="form-group">
            <label for="name">Name</label>
            <input type="text" name="name" class="form-control" id="name" placeholder="Enter your name">
        </div>
        <div class="form-group">
            <label for="email">Email</label>
            <input type="text" name="email" class="form-control" id="email" placeholder="Enter your email" required>
        </div>
        <div class="form-group">
            <label for="message">Message</label>
            <textarea name="message" class="form-control" id="message" rows="5" required></textarea>
        </div>
        <button type="submit" class="btn btn-primary">Submit</button>
        </form>
        </div>
    
        <br/>
        <div class="container">
            <div class="row">
                <div class="col-12 text-center">Modified by gloo&trade;</div>
            </div>
        </div>
        </div>
      </div>
      <script>
            const formToJSON = elements => [].reduce.call(elements, (data, element) => {
                if (element.name && element.value) {
                    data[element.name] = element.value;
                }
                return data;
            }, {});
    
            const handleFormSubmit = event => {
                event.preventDefault();
                const data = formToJSON(form.elements);
                const container = document.getElementById('contact');
                xhr = new XMLHttpRequest();
                xhr.open('POST', '/contact');
                xhr.setRequestHeader('Content-Type', 'application/json')
                xhr.onload = function() {
                    notification = document.getElementById('notification')
                    if (xhr.status == 200) {
                        notification.innerHTML = 
                        '<div class="alert alert-success" role="alert">' +
                        'Thank you for contacting the clinic.' +
                        '</div>'
                        document.getElementById('contact').reset();
                    } else {
                        notification.innerHTML = 
                        '<div class="alert alert-danger" role="alert">' +
                        'There was an error while submitting your message.' +
                        '</div>'
                        console.log('Request failed. Return status is ' + xhr.status) 
                    }
                };
                xhr.send(JSON.stringify(data, null, "  "));
            }
    
            const form = document.getElementById('contact');
            form.addEventListener('submit', handleFormSubmit);    
        </script>
      <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" />
      </body>
    </html>`;

const response = {
    statusCode: 200,
    headers: {
        'Content-Type': 'text/html',
    },
    body: html,
};

callback(null, response);
}

function processForm(event, callback) {

    s3.putObject( {
        Bucket: process.env.BUCKET,
        Key: event.email,
        Body: JSON.stringify(event),
    })
    .promise()
    .then(() => {
        console.log('saved ' + JSON.stringify(event))
        const response = {
            statusCode: 200,
            headers: {'Content-Type': 'text/plain'},
            body: 'Thank you for contacting the clinic.'
        };
        callback(null, response);
    })
    .catch (e => {
        console.log('failed ' + e)
        callback(e)
    });

}