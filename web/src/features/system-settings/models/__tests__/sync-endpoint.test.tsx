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
import { render, screen, within } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { expect, test } from 'vitest'

import type { UpstreamChannel } from '../../types'
import { ChannelSelectorDialog } from '../channel-selector-dialog'

const channel = {
  id: 4,
  name: 'Reseller',
  base_url: 'https://api.example.com',
  status: 1,
} as UpstreamChannel

function Fixture() {
  const [endpoints, setEndpoints] = useState<Record<number, string>>({
    4: '/api/pricing',
  })
  return (
    <ChannelSelectorDialog
      open
      onOpenChange={() => {}}
      channels={[channel]}
      selectedChannelIds={[]}
      onSelectedChannelIdsChange={() => {}}
      channelEndpoints={endpoints}
      onChannelEndpointsChange={setEndpoints}
      onConfirm={() => {}}
    />
  )
}

test('choosing the custom sync endpoint shows a field for the address and keeps what is typed', async () => {
  const user = userEvent.setup()
  render(<Fixture />)
  const row = screen.getByRole('row', { name: /Reseller/ })
  await user.click(within(row).getByRole('combobox'))
  await user.click(await screen.findByRole('option', { name: 'custom' }))
  const field = await within(row).findByPlaceholderText('/your/endpoint')
  await user.type(field, 'https://example.com/api/pricing')
  expect(field).toHaveValue('https://example.com/api/pricing')
  expect(within(row).getByRole('combobox')).toHaveTextContent('custom')
})
