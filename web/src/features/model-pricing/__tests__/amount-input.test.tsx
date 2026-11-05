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
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { useState } from 'react'
import { expect, test } from 'vitest'

import { PricingAmountInput } from '../pricing-amount-input'

function Fixture() {
  const [usd, setUsd] = useState('')
  return (
    <>
      <PricingAmountInput
        aria-label='Input price'
        value={usd}
        onChange={setUsd}
      />
      <output aria-label='Saved price'>{usd}</output>
    </>
  )
}

test('a decimal comma from a comma-locale keypad is accepted as the separator', async () => {
  const user = userEvent.setup()
  render(<Fixture />)
  const input = screen.getByRole('textbox', { name: 'Input price' })
  await user.type(input, '2,3077')
  expect(input).toHaveValue('2.3077')
  expect(screen.getByRole('status', { name: 'Saved price' })).toHaveTextContent(
    '2.3077'
  )
})
