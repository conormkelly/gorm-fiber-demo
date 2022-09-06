const newman = require('newman')
const { execSync, spawn } = require('child_process')

// CONFIG
process.env.APP_PORT = '3000'
process.env.APP_DB_TYPE = 'MYSQL'
process.env.APP_DB_CONN_STRING = 'TBD'

const APP_EXECUTABLE =
  process.platform === 'win32' ? 'fiber-demo.exe' : 'fiber-demo'

const STARTUP_TIMEOUT_MS = 30000

async function main() {
  console.time('Time Taken')

  // Build the app
  console.log('RUNNING: BUILD')
  execSync('go build', { encoding: 'utf-8' })

  // Create a subprocess of the app
  const app = spawn('./' + APP_EXECUTABLE)

  // Set a timer in case app doesn't start successfully
  const startupTimer = setTimeout(() => {
    console.error(
      `\nFAILED: app startup took greater than ${STARTUP_TIMEOUT_MS} ms`
    )
    app.kill()
  }, STARTUP_TIMEOUT_MS)

  // Wait for Fiber to start
  app.stdout.on('data', async (chunk) => {
    const output = chunk.toString('utf-8')
    if (
      output.includes(`bound on host 0.0.0.0 and port ${process.env.APP_PORT}`)
    ) {
      clearTimeout(startupTimer)
      console.log('APP: READY')
      await executePostmanCollection('./Fiber_Gorm.postman_collection.json')
      app.kill()
    }
  })

  app.on('error', (err) => {
    console.log('ERROR:', err)
    process.exit(1)
  })

  app.on('exit', (code, signal) => {
    console.log(`App exited with code ${code} and signal ${signal}`)
    console.timeEnd('Time Taken')
    clearTimeout(startupTimer)
    process.exit(1)
  })
}

function executePostmanCollection(path) {
  return new Promise((resolve, reject) => {
    newman.run(
      {
        collection: require(path),
        reporters: 'cli',
      },
      function (err) {
        if (err) {
          reject(err)
        } else {
          resolve()
        }
      }
    )
  })
}

main()
