// SPDX-FileCopyrightText: 2025 Forkbomb BV
//
// SPDX-License-Identifier: AGPL-3.0-or-later

import { Command } from 'commander';
import { runFromFile } from '@stepci/runner';

const program = new Command();

program
  .option('-p, --path <file>', 'Path to the test file', './stepci.yaml')
  .option('-s, --secret <key=value...>', 'Secrets as key=value pairs')
  .option('-e, --env <key=value...>', ' Environment variables  as key=value pairs');

program.parse(process.argv);
const options = program.opts();

// Parse secrets into an object
const secrets = {};
if (options.secret) {
  options.secret.forEach((pair) => {
    const match = pair.match(/^([^=]+)=(.*)$/); // Match key=value pattern

    if (match) {
      const key = match[1].trim();
      let value = match[2].trim();
      if (value.startsWith('"') && value.endsWith('"')) {
        value = value.slice(1, -1);
      }
      secrets[key] = value;
    }
  });
}

const env = {};
if (options.env) {
  options.envforEach((pair) => {
    const match = pair.match(/^([^=]+)=(.*)$/); // Match key=value pattern

    if (match) {
      const key = match[1].trim();
      let value = match[2].trim();
      if (value.startsWith('"') && value.endsWith('"')) {
        value = value.slice(1, -1);
      }
      env[key] = value;
    }
  });
}

const workflowOptions = {
    secrets: secrets,
    env: env
  };

runFromFile(options.path, workflowOptions).then((result) => {
  const { passed, tests } = result.result;
  if (passed) {
    const lastTest = tests[tests.length - 1];
    const lastStep = lastTest.steps[lastTest.steps.length - 1];

    if (lastStep?.captures) {
      console.log(JSON.stringify(lastStep.captures, null, 2));
    } else {
      console.log("No captures found in the last step.");
    }
  } else {
    console.error("❌ Workflow failed. Details:\n");

    tests.forEach((test) => {
      test.steps.forEach((step) => {
        if (!step.passed) {
          console.error(`🔴 Step Failed: ${step.name}`);
          console.error(`  🌍 URL: ${step.request?.url}`);
          console.error(`  📡 Method: ${step.request?.method}`);

          if (step.checks) {
            let hasErrors = false;
            let errorMessages = "  ❌ Failed Checks:\n";

            for (const key in step.checks) {
              const check = step.checks[key];

              if ('expected' in check && 'given' in check) {
                if (JSON.stringify(check.expected) !== JSON.stringify(check.given)) {
                  hasErrors = true;
                  errorMessages += `    - ${key}: Expected ${JSON.stringify(check.expected)}, but got ${JSON.stringify(check.given)}\n`;
                }
              } else {
                for (const subKey in check) {
                  const subCheck = check[subKey];
                  if (JSON.stringify(subCheck.expected) !== JSON.stringify(subCheck.given)) {
                    hasErrors = true;
                    errorMessages += `    - ${key}.${subKey}: Expected ${JSON.stringify(subCheck.expected)}, but got ${JSON.stringify(subCheck.given)}\n`;
                  }
                }
              }
            }

            if (hasErrors) {
              console.error(errorMessages);
            }
          }

          if (step.response) {
            console.error("  📩 Response:");
            console.error(`    - Status: ${step.response.status} ${step.response.statusText}`);
            console.error(`    - Headers: ${JSON.stringify(step.response.headers, null, 2)}`);

            if (step.response.body) {
              const responseBody = Buffer.from(step.response.body).toString('utf-8');
              try {
                const responseJson = JSON.parse(responseBody);
                console.error(`    - Body (JSON):\n${JSON.stringify(responseJson, null, 2)}`);
              } catch {
                console.error(`    - Body (Raw Text):\n${responseBody}`);
              }
            }
          }

          console.error("\n");
        }
      });
    });

    process.exit(1);
  }
});