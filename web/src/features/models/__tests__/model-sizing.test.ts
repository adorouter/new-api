/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { expect, test } from 'vitest'

import {
  transformFormDataToModelPayload,
  transformModelToFormDefaults,
  writeModelSizing,
} from '../lib/model-form'

const stored = {
  id: 7,
  model_name: 'example-model',
  status: 1,
  sync_official: 1,
  name_rule: 0,
  created_time: 1,
  updated_time: 1,
  metadata:
    '{"contextWindow":1000000,"maxOutputTokens":64000,"isReasoning":true}',
}

test('saving the metadata form sends the stored metadata back unchanged', () => {
  const form = transformModelToFormDefaults(stored)
  expect(form.context_window).toBe('1000000')
  expect(form.max_output_tokens).toBe('64000')
  const payload = transformFormDataToModelPayload({
    ...form,
    description: 'Edited',
  })
  expect(JSON.parse(payload.metadata ?? '')).toEqual(
    JSON.parse(stored.metadata)
  )
})

test('editing the sizing fields keeps every other metadata key', () => {
  const metadata = writeModelSizing(stored.metadata, '200000', '')
  expect(JSON.parse(metadata)).toEqual({
    contextWindow: 200000,
    isReasoning: true,
  })
})

test('sizing on a model without metadata creates only the sizing keys', () => {
  expect(writeModelSizing('', '128000', '8192')).toBe(
    '{"contextWindow":128000,"maxOutputTokens":8192}'
  )
  expect(writeModelSizing('', '', '')).toBe('')
})

test('metadata that is not a JSON object is left as stored', () => {
  expect(writeModelSizing('not json', '128000', '')).toBe('not json')
})
