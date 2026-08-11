

export let config = {
    offsetPost: 0,
    offsetMessages:0
}

// A post carries this reaction id when the signed in user did not react yet:
// the backend COALESCEs the missing reactions row to the nil uuid.
export const ZERO_UUID = "00000000-0000-0000-0000-000000000000"