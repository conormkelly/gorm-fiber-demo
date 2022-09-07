const newman = require('newman')
const { execSync, spawn } = require('child_process')

// Config
const SCRIPT_TIMEOUT_MS = 30000

async function main() {
  console.time('Time taken')

  // Create a subprocess of the docker-compose process
  console.log('RUNNING: docker-compose up')
  const dockerProcess = spawn('docker-compose', ['up'])
  dockerProcess.stdout.setEncoding('utf8')

  // Set a timer in case the app doesn't start successfully
  const startupTimer = setTimeout(() => {
    console.error(
      `\nFAILED: app startup took greater than ${SCRIPT_TIMEOUT_MS} ms`
    )
    shutdownDockerProcess(dockerProcess)
  }, SCRIPT_TIMEOUT_MS)

  // Wait for Fiber to start within docker subprocess
  // by monitoring the stdout
  dockerProcess.stdout.on('data', async (data) => {
    // const output = data.toString('utf-8')
    console.log(data)
    if (data.includes('bound on')) {
      // App has started up, so clear the timer
      clearTimeout(startupTimer)

      let exitCode = 0
      try {
        await executePostmanCollection('./Fiber_Gorm.postman_collection.json')        
      } catch (err) {
        exitCode = 1
        console.error('An error occured executing postman collection:', err.message)
      } finally {
        shutdownDockerProcess(dockerProcess)
      }
      process.exit(exitCode)
    }
  })

  dockerProcess.on('error', (err) => {
    console.log('ERROR:', err)
    process.exit(1)
  })

  dockerProcess.on('exit', (code, signal) => {
    console.log(`App exited with code ${code} and signal ${signal}`)
    console.timeEnd('Time taken')
    clearTimeout(startupTimer)
    process.exit(1)
  })
}

function shutdownDockerProcess(process) {
  try {
    console.log('RUNNING: docker-compose down')
    execSync('docker-compose down')
  } catch (err) {
    console.log('ERROR: Failed to run docker-compose down')
    process.kill()
  }
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
