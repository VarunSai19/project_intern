'use strict';
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const bodyParser = require('body-parser');
const http = require('http')
const util = require('util');
var SHA256 = require("crypto-js/sha256");
const mongoose = require('mongoose')
const express = require('express')
const app = express();
const dbURI = 'mongodb+srv://varun:varun1234@telco-project.rf8w0.mongodb.net/Telco-project?retryWrites=true&w=majority'

const cors = require('cors');
const constants = require('./config/constants.json')

const host = process.env.HOST || constants.host;
const port = process.env.PORT || constants.port;

const helper = require('./app/helper')
var user_hash_dict = require('./app/helper').user_hash_dict
const invoke = require('./app/invoke')
const qscc = require('./app/qscc')
const query = require('./app/query')
// const PasswordHash = require('./models/schema_pass');
const Data = require('./schema_data');
const channelName = "mychannel"
const chaincodeName = "fabcar"

mongoose.connect(dbURI,{useNewUrlParser:true,useUnifiedTopology:true})
    .then((result) => {
        var server = http.createServer(app).listen(port, function () { console.log(`Server started on ${port}`) });
        logger.info('****************** SERVER STARTED ************************');
        logger.info('***************  http://%s:%s  ******************', host, port);
        server.timeout = 240000;
    })
    .catch((err) => console.log(err));


app.use(express.static('public'));
app.use("/css",express.static(__dirname+'public/css'))
// app.use("/img",express.static(__dirname+'public/img'))

app.set('views','./views');
app.set('view engine', 'ejs');

app.options('*', cors());
app.use(cors());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
    extended: false
}));


logger.level = 'debug';

var aadhardata = {}
aadhardata["1"] = {"Name":"varun","Gender":"M","DateOfBirth":"27011999","Address":"Nellore"};
aadhardata["12"] = {"Name":"varun","Gender":"M","DateOfBirth":"27011999","Address":"Nellore"};
aadhardata["123"] = {"Name":"varun","Gender":"M","DateOfBirth":"27011999","Address":"Nellore"};
aadhardata["1234"] = {"Name":"varun","Gender":"M","DateOfBirth":"27011999","Address":"Nellore"};
aadhardata["12345"] = {"Name":"varun","Gender":"M","DateOfBirth":"27011999","Address":"Nellore"};

// app.use((req, res, next) => {
//     logger.debug('New req for %s', req.originalUrl);
//     if (req.originalUrl.indexOf('/user/varun/add_data') >= 0 || req.originalUrl.indexOf('/users') >= 0 || req.originalUrl.indexOf('login') >= 0 || req.originalUrl.indexOf('/register') >= 0) {
//         return next();
//     }
//     return next();
// });

// var server = http.createServer(app).listen(port, function () { console.log(`Server started on ${port}`) });
// logger.info('****************** SERVER STARTED ************************');
// logger.info('***************  http://%s:%s  ******************', host, port);
// server.timeout = 240000;


function getErrorMessage(field) {
    var response = {
        success: false,
        message: field + ' field is missing or Invalid in the request'
    };
    return response;
}

app.get('/', async function(req,res){
    res.render('index',{title:'Home'})
});

app.get('/register',async function (req, res) {
    res.render('Login',{title:"Register"})
});

// Register and enroll user
app.post('/register', async function (req, res) {
    var username = req.body.username;
    var password = req.body.password;
    var usertype = req.body.usertype;
    var orgName = helper.getOrg(usertype)
    logger.debug('End point : /register');
    logger.debug('User name : ' + username);
    logger.debug('Password  : ' + password);
    logger.debug('Usertype  : ' + usertype);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!password) {
        res.json(getErrorMessage('\'password\''));
        return;
    }
    if (!usertype) {
        res.json(getErrorMessage('\'usertype\''));
        return;
    }

    let response = await helper.Register(username, password,usertype);

    logger.debug('-- returned from registering the username %s for organization %s', username, orgName);
    if (response && typeof response !== 'string') {
        logger.debug('Successfully registered the username %s for organization %s', username, orgName);
        // res.json(response);
        var pass_hash = SHA256(username+password)
        pass_hash = JSON.stringify(pass_hash["words"]);
        console.log(pass_hash);
        // const pw_data = new PasswordHash({
        //     username:username,
        //     password_hash:pass_hash
        // });
        // pw_data.save();
        res.render('success',{username,title:"success"})
    } else {
        logger.debug('Failed to register the username %s for organization %s with::%s', username, orgName, response);
        // res.json({ success: false, message: response });
        res.render('failure',{username,title:"failed"})
    }

});

// Login 
// Here we need to develop front end page for each login files
app.post('/Adminlogin', async function (req, res) {
    res.render('Login',{title:"Admin Login"})
});

app.post('/Adminlogin', async function (req, res) {
    var username = req.body.username;
    const user_present = helper.isUserRegistered(username,"Org1")
    if(!user_present) 
    {
        console.log(`An identity for the user ${username} not exists`);
        var response = {
            success: false,
            message: username + ' was not enrolled',
        };
        return response
    }
    var password = req.body.password;
    var usertype = req.body.usertype;
    var orgName = helper.getOrg(usertype);
    logger.debug('End point : /login');
    logger.debug('User name : ' + username);
    logger.debug('Password  : ' + password);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!password) {
        res.json(getErrorMessage('\'Password\''));
        return;
    }
    // This should be changed as per the file in GCP
    var pass_hash = SHA256(username+password)
    if(JSON.stringify(user_hash_dict[username]["password_hash"]["words"]) !== JSON.stringify(pass_hash["words"]))
    {
        res.json({success: false, message: "Invalid Credentials" });
    }
    res.json({ success: true });
});


app.post('/Userlogin', async function (req, res) {
    var username = req.body.username;

    const user_present = helper.isUserRegistered(username,"Org2")
    if(!user_present) 
    {
        console.log(`An identity for the user ${username} not exists`);
        var response = {
            success: false,
            message: username + ' was not enrolled',
        };
        return response
    }
    var password = req.body.password;
    var usertype = "customer";
    var orgName = helper.getOrg(usertype);
    logger.debug('End point : /login');
    logger.debug('User name : ' + username);
    logger.debug('Password  : ' + password);
    if (!username) {
        res.json(getErrorMessage('\'username\''));
        return;
    }
    if (!password) {
        res.json(getErrorMessage('\'Password\''));
        return;
    }
    // This should be changed as per the file in GCP
    var pass_hash = SHA256(username+password)
    if(JSON.stringify(user_hash_dict[username]["password_hash"]["words"]) !== JSON.stringify(pass_hash["words"]))
    {
        var url_new = '/user'+username
        res.json({success: false, message: "Invalid Credentials" });
    }
    res.json({ success: true });
});


app.get('/dealer/getSimCard',async function(req,res){
    res.render('dealer_page',{title:"New User",phonenumbers})
});

app.post('/dealer/getSimCard' ,async function (req,res){
    try{
        var password = req.body.password;

        var args = {};
        args["AadharNumber"] = JSON.stringify(req.body.AadharNumber);
        args["Address"] = req.body.Address;
        args["DateOfBirth"] = req.body.DateOfBirth;
        args["Name"] = req.body.Name;
        args["AltenativeNumber"] = JSON.stringify(req.body.AltenativeNumber);
        // args["PhoneNumber"] = JSON.stringify(req.body.PhoneNumber);
        args["Gender"] = req.body.Gender;
        
        console.log(`Input is ${args}`)

        var actual_data = aadhardata[args["AadharNumber"]]

        console.log(`Actual data is ${actual_data}`)

        if(!actual_data)
        {
            result = "Aadhar data doesnt exist in server";
            const response_payload = {
                result: result,
                error: error.name,
                errorData: error.message
            }
            res.send(response_payload)
        }
        console.log("Aadhar data present in the system.")
        const valid = await helper.ValidateAadhar(actual_data,args);
        if(!valid)
        {
            result = "Data provided is not matched by actual data.";
            const response_payload = {
                result: result,
                error: error.name,
                errorData: error.message
            }
            res.send(response_payload)
        }
        console.log("Aadhar data matched.")
        let response = await helper.Register(args["PhoneNumber"], password,"customer");
        
        console.log("User created...")

        let message = await invoke.invokeTransaction("CreateData",args["PhoneNumber"],args);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    }
    catch(error)
    {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});



app.post('/user/:userid/ChangeDetails' ,async function (req,res){
    try{
        var password = req.body.password;
        var username = req.params.userid;
        var pass_hash = SHA256(username+password)
        if(JSON.stringify(user_hash_dict[username]["password_hash"]["words"]) !== JSON.stringify(pass_hash["words"]))
        {
            res.json({success: false, message: "Invalid Credentials" });
            return;
        }
        var args = {};
        args["AadharNumber"] = JSON.stringify(req.body.AadharNumber);
        args["Address"] = req.body.Address;
        args["DateOfBirth"] = req.body.DateOfBirth;
        args["Name"] = req.body.Name;
        args["Gender"] = req.body.Gender;
        
        console.log(`Input is ${args}`)

        var actual_data = aadhardata[args["AadharNumber"]]

        console.log(`Actual data is ${actual_data}`)

        if(!actual_data)
        {
            result = "Aadhar data doesnt exist in server";
            const response_payload = {
                result: result,
                error: error.name,
                errorData: error.message
            }
            res.send(response_payload)
        }
        console.log("Aadhar data present in the system.")
        const valid = await helper.ValidateAadhar(actual_data,args);
        if(!valid)
        {
            result = "Data provided is not matched by actual data.";
            const response_payload = {
                result: result,
                error: error.name,
                errorData: error.message
            }
            res.send(response_payload)
        }
        console.log("Aadhar data matched.")

        let message = await invoke.invokeTransaction("ChangeData",args["PhoneNumber"],args);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    }
    catch(error)
    {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.post('/user/:userid/services/BuyService' ,async function (req,res){
    try{
        var password = req.body.password;
        var username = req.params.userid;
        var pass_hash = SHA256(username+password)
        if(JSON.stringify(user_hash_dict[username]["password_hash"]["words"]) !== JSON.stringify(pass_hash["words"]))
        {
            res.json({success: false, message: "Invalid Credentials" });
            return;
        }
        var args = {};
        args["Service_name"] = req.body.Service_name;
        args["Price"] = req.body.Price;
        
        console.log(`Input is ${args}`)

        let message = await invoke.invokeTransaction("BuyService",username,args);
        console.log(`message result is : ${message}`)

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);

    }
    catch(error)
    {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});


app.get('/admin/:username/GetIdentity', async function (req, res) {
    try{
        let username = req.params.username
        let message = await query.query(null, "GetSubmittingClientIdentity",username,"Org1");
        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    }
    catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});


app.get('/admin/:username/GetCustomerByPhoneNumber', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');
        
        let args = req.query.args;
        let username = req.params.username
        logger.debug('args : ' + args);

        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }
        console.log('args==========', args);
        args = args.replace(/'/g, '"');
        args = JSON.parse(args);
        logger.debug(args);

        let message = await query.query(args, "GetDataByPhoneNumber",username,"Org1");

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});



app.get('/admin/:username/GetAllCustomers', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');
        
        let username = req.params.username
        let message = await query.query(null,"QueryAllData",username,"Org1");

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }
        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});


app.get('/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        // var channelName = req.params.channelName;
        // var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`)
        let args = req.query.args;
        let fcn = req.query.fcn;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn : ' + fcn);
        logger.debug('args : ' + args);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }
        console.log('args==========', args);
        args = args.replace(/'/g, '"');
        args = JSON.parse(args);
        logger.debug(args);

        let message = await query.query(channelName, chaincodeName, args, fcn, req.username, req.orgname);

        const response_payload = {
            result: message,
            error: null,
            errorData: null
        }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});

app.get('/qscc/channels/:channelName/chaincodes/:chaincodeName', async function (req, res) {
    try {
        logger.debug('==================== QUERY BY CHAINCODE ==================');

        var channelName = req.params.channelName;
        var chaincodeName = req.params.chaincodeName;
        console.log(`chaincode name is :${chaincodeName}`)
        let args = req.query.args;
        let fcn = req.query.fcn;
        // let peer = req.query.peer;

        logger.debug('channelName : ' + channelName);
        logger.debug('chaincodeName : ' + chaincodeName);
        logger.debug('fcn : ' + fcn);
        logger.debug('args : ' + args);

        if (!chaincodeName) {
            res.json(getErrorMessage('\'chaincodeName\''));
            return;
        }
        if (!channelName) {
            res.json(getErrorMessage('\'channelName\''));
            return;
        }
        if (!fcn) {
            res.json(getErrorMessage('\'fcn\''));
            return;
        }
        if (!args) {
            res.json(getErrorMessage('\'args\''));
            return;
        }
        console.log('args==========', args);
        args = args.replace(/'/g, '"');
        args = JSON.parse(args);
        logger.debug(args);

        let response_payload = await qscc.qscc(channelName, chaincodeName, args, fcn, req.username, req.orgname);

        // const response_payload = {
        //     result: message,
        //     error: null,
        //     errorData: null
        // }

        res.send(response_payload);
    } catch (error) {
        const response_payload = {
            result: null,
            error: error.name,
            errorData: error.message
        }
        res.send(response_payload)
    }
});