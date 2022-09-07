const newman = require('newman')
const { execSync, spawn } = require('child_process')

const SCRIPT_TIMEOUT_MS = 30000

async function main() {
  console.time('Time taken')

  // Create a subprocess of the docker-compose process
  console.log('RUNNING: docker-compose-up')
  const dockerProcess = spawn('docker-compose', ['up'])

  // Set a timer in case app doesn't start successfully
  const startupTimer = setTimeout(() => {
    console.error(
      `\nFAILED: app startup took greater than ${SCRIPT_TIMEOUT_MS} ms`
    )
    try {
      execSync('docker-compose down')
    } catch (err) {
      console.log('ERROR: Failed to run docker-compose down')
      dockerProcess.kill()
    }
  }, SCRIPT_TIMEOUT_MS)

  // Wait for Fiber to start
  dockerProcess.stdout.on('data', async (chunk) => {
    const output = chunk.toString('utf-8')
    console.log(output)
    if (output.includes('bound on')) {
      clearTimeout(startupTimer)
      console.log('APP: READY')
      await executePostmanCollection('./Fiber_Gorm.postman_collection.json')
      try {
        execSync('docker-compose down')
      } catch (err) {
        console.log('ERROR: Failed to run docker-compose down')
        dockerProcess.kill()        
      }
      process.exit(0)
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
