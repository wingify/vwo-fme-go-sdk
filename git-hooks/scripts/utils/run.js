/**
 * Copyright 2024 Wingify Software Pvt. Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

const childProcess = require('child_process');
const AnsiColorEnum = require('../enums/AnsiColorEnum');

function _exec(cmd) {
  return childProcess
    .execSync(cmd, {stdio: 'inherit'})
}

function run({
  name = 'Unknwon',
  command = `echo 'No command was specified!'`,
  successMessage = 'Passed',
  failureMessage = 'Failed'
}) {
  console.log(`${AnsiColorEnum.CYAN}\nRunning ${name}${AnsiColorEnum.RESET} : ${AnsiColorEnum.YELLOW}${command}${AnsiColorEnum.RESET}\n`);
  try {
    _exec(command);
    console.log(`${AnsiColorEnum.GREEN}\n\n${name} ${successMessage}${AnsiColorEnum.RESET}\n\n`);
  } catch (e) {
    console.log(`${AnsiColorEnum.RED}\n\n${name} ${failureMessage}${AnsiColorEnum.RESET}\n\n`, e);
    process.exit(1);
  }
}

module.exports = {
  run
};