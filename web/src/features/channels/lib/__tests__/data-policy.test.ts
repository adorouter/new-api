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

import { channelSchema } from '../../types'
import {
  CHANNEL_FORM_DEFAULT_VALUES,
  buildSettingJSON,
  channelFormSchema,
  transformChannelToFormDefaults,
} from '../channel-form'

function storedChannel(setting: string) {
  return channelSchema.parse({
    id: 7,
    name: 'Reseller',
    type: 1,
    key: '',
    models: 'example-model',
    group: 'default',
    status: 1,
    created_time: 1,
    test_time: 0,
    response_time: 0,
    balance_updated_time: 0,
    other: '',
    remark: '',
    setting,
  })
}

test('saving a channel keeps its data policy and settings the form does not manage', () => {
  const setting = JSON.stringify({
    proxy: '',
    data_policy: {
      provider: 'Example AI',
      seat_country: 'SG',
      prompts: 'trains',
    },
    upstream_cost_markup: 1.2,
  })
  const form = channelFormSchema.parse(
    transformChannelToFormDefaults(storedChannel(setting))
  )
  expect(form.data_policy_provider).toBe('Example AI')
  expect(form.data_policy_seat_country).toBe('SG')
  expect(form.data_policy_prompts).toBe('trains')
  expect(form.data_policy_transfer).toBe('unknown')

  const saved = JSON.parse(buildSettingJSON(form))
  expect(saved.data_policy).toEqual({
    provider: 'Example AI',
    seat_country: 'SG',
    prompts: 'trains',
  })
  expect(saved.upstream_cost_markup).toBe(1.2)
})

test('clearing every data policy field removes the data policy from the setting', () => {
  const setting = JSON.stringify({
    data_policy: { provider: 'Example AI', prompts: 'none' },
  })
  const form = channelFormSchema.parse({
    ...transformChannelToFormDefaults(storedChannel(setting)),
    data_policy_provider: '',
    data_policy_prompts: 'unknown',
  })
  expect(JSON.parse(buildSettingJSON(form))).not.toHaveProperty('data_policy')
})

test('the provider country must be a two-letter code', () => {
  const result = channelFormSchema.safeParse({
    ...CHANNEL_FORM_DEFAULT_VALUES,
    name: 'Reseller',
    type: 1,
    key: 'test-key',
    models: 'example-model',
    data_policy_seat_country: 'USA',
  })
  expect(result.success).toBe(false)
})
