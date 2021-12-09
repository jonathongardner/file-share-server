const express = require('express');
const fileUpload = require('express-fileupload');
const bodyParser = require('body-parser');
const path = require('path');
const fs = require('fs');

const app = express();
app.use(fileUpload({ createParentPath: true }));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: true }));

const port = process.env.PORT || 3000;
const uploadPath = process.env.UPLOAD_PATH || path.join(__dirname, 'tmp')

const renderJson = (res, json, code = 200) => {
  res.setHeader('Content-Type', 'application/json');
  res.status(code);
  res.send(JSON.stringify(json));
}
const renderError = (res, message, code = 200) => {
  res.status(code);
  res.send(`<div>${message}</div>`);
}

const checkPath = (fullPath, { file = true, folder = true }) => {
  if (!file && !folder) {
    if (fs.existsSync(fullPath)) {
      throw { message: 'Path already exists', statusCode: 422 }
    }
    return { file: false, directory: false, fullPath: fullPath }
  }

  const sync = fs.lstatSync(fullPath)
  if (file && sync.isFile()) {
    return { file: true, directory: false, fullPath: fullPath }
  }
  if (folder && sync.isDirectory()) {
    return { file: false, directory: true, fullPath: fullPath }
  }
  throw { message: 'Path not found', statusCode: 404 }
}

const errorHandler = (callback, options = {}) => {
  return (req, res) => {
    console.log('api: ' + req.path)
    try {
      const fullPath = path.join(uploadPath, req.params[0] || '')
      callback(req, res, checkPath(fullPath, options))
    } catch (e) {
      console.log(e)
      if(e.statusCode) {
        renderJson(res, { error: [e.message] }, e.statusCode)
      } else {
        renderJson(res, { error: ['Unknown error'] }, 500)
      }
    }
  }
}

//----------Routes-----------
// ^\/api\/files\/((?:(?!\.\.).)*)$/ if need to not allow .. but think node handles
const apiRoute = /^\/api\/files(?:\/(.*))?$/

// download file/list folder
app.get(apiRoute, errorHandler((req, res, { file, fullPath }) => {
  if (file) {
    console.log('Downloading file ' + fullPath)
    res.download(fullPath);
  } else {
    console.log('Listing folder ' + fullPath)
    renderJson(res, {
      data: fs.readdirSync(fullPath, { withFileTypes: true }).map(v => {
        // const file = path.join(fullPath, v.name)
        // const stats = fs.statSync(file)
        return { name: v.name, file: v.isFile(), directory: v.isDirectory() }
      }).filter(v => v.file || v.directory)
    })
  }
}));

// delete file/folder
app.delete(apiRoute, errorHandler((req, res, { file, fullPath }) => {
  if (fullPath === uploadPath) {
    return renderJson(res, { error: 'Cant delete base folder' }, 422);
  }
  console.log('Deleting ' + fullPath)
  if (file) {
    fs.unlinkSync(fullPath)
  } else {
    fs.rmdirSync(fullPath, { recursive: true });
  }
  renderJson(res, {});
}))

// create file/folder
app.post(apiRoute, errorHandler(async (req, res, { fullPath }) => {
  if (req.files && req.files.file) {
    console.log('Creating file ' + fullPath)
    await req.files.file.mv(fullPath)
  } else {
    console.log('Creating directory ' + fullPath)
    fs.mkdirSync(fullPath, { recursive: true })
  }
  renderJson(res, {});
}, { file: false, folder: false }));
//----------Routes-----------

//----------UI-----------
app.get('/*', (req, res) => {

  console.log('UI: ' + req.path)
  try {
    if (req.path.startsWith('/static')) {
      res.sendFile(path.join(__dirname, req.path));
    } else {
      const fullPath = path.join(uploadPath, req.path)
      checkPath(fullPath, { file: false }) // only want to see if directory at this place
      res.sendFile(path.join(__dirname, 'static/index.html'));
    }
  } catch (e) {
    console.log(e)
    if(e.statusCode) {
      renderError(res, e.message, e.statusCode)
    } else {
      renderError(res, 'Unknown error', 500)
    }
  }
});
//----------UI-----------

app.listen(port, function () {
  console.log(`File Server listening on port ${port}!`);
  console.log(`Folder ${uploadPath}!`);
});
