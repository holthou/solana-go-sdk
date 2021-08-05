package types

import (
	"encoding/hex"
	"fmt"
	"github.com/mr-tron/base58"
	"reflect"
	"testing"

	"github.com/portto/solana-go-sdk/common"
)

func TestTransaction_Serialize(t *testing.T) {
	type fields struct {
		Signatures []Signature
		Message    Message
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			fields: fields{
				Signatures: []Signature{[]byte{189, 98, 67, 19, 102, 99, 124, 234, 70, 209, 28, 10, 33, 66, 167, 162, 222, 122, 16, 68, 248, 129, 46, 111, 221, 255, 40, 40, 236, 84, 233, 213, 234, 185, 235, 222, 155, 204, 139, 164, 184, 155, 32, 54, 151, 73, 235, 65, 200, 76, 127, 111, 244, 72, 183, 208, 21, 247, 114, 176, 181, 21, 77, 8}},
				Message: Message{
					Header: MessageHeader{
						NumRequireSignatures:        1,
						NumReadonlySignedAccounts:   0,
						NumReadonlyUnsignedAccounts: 1,
					},
					Accounts: []common.PublicKey{
						common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
						common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"),
						common.SystemProgramID,
					},
					RecentBlockHash: "FwRYtTPRk5N4wUeP87rTw9kQVSwigB6kbikGzzeCMrW5",
					Instructions: []CompiledInstruction{
						{
							ProgramIDIndex: 2,
							Accounts:       []int{0, 1},
							Data:           []byte{2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
						},
					},
				},
			},
			want:    []byte{1, 189, 98, 67, 19, 102, 99, 124, 234, 70, 209, 28, 10, 33, 66, 167, 162, 222, 122, 16, 68, 248, 129, 46, 111, 221, 255, 40, 40, 236, 84, 233, 213, 234, 185, 235, 222, 155, 204, 139, 164, 184, 155, 32, 54, 151, 73, 235, 65, 200, 76, 127, 111, 244, 72, 183, 208, 21, 247, 114, 176, 181, 21, 77, 8, 1, 0, 1, 3, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 221, 244, 189, 59, 8, 252, 7, 91, 129, 169, 22, 151, 32, 104, 208, 131, 64, 75, 232, 201, 77, 13, 187, 220, 103, 232, 190, 100, 35, 210, 17, 42, 1, 2, 2, 0, 1, 12, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				Signatures: tt.fields.Signatures,
				Message:    tt.fields.Message,
			}
			got, err := tx.Serialize()
			if (err != nil) != tt.wantErr {
				t.Errorf("Transaction.Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transaction.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

//反解交易
func TestTransaction_Serialize11(t *testing.T) {
	encodedStr := "02b3302b6fd89a72ab7e5e28ef91783c0d7a776c29ccae00c3ef651e3d5bf17521e78888317bb0a5663e7beb864de1ee729d56c188691b19b6613efa9541c07d041c212efdfa0b70dfb3e863cd8f607a773c56c52f73ac901af34c4a1f6c058f785df85d15b85754dfd1af5f027b3bccc76abeadfe01f9df2aeed6b8a9fd19140002010105f091a7f1d3a414f944a9cecf6bd0dab1cda1fec5c809860341d2403897f4285ed8ae77a0477191544e34285e5f1986232bfb1cef53f5b0138f7db33a432c221f26193c712a6f06efd4b732a46724b017252e0f60340c32e1c492f92af1143c577c201260fc8a7bd77e5c146bd61c1d55914156850ba81216224c2868f0e160be06ddf6e1d765a193d9cbe146ceeb79ac1cb485ed5f5b37913a8cf5857eff00a9f02f202f062feedd0d4b9c5e7ca76e95de561362f66aebd6e1fbda42ce6c3a370104030302010903c234000000000000"
	test, _ := hex.DecodeString(encodedStr)
	fmt.Println(base58.Encode(test))
	trx := MustTransactionDeserialize(test)

	for i, sg := range trx.Signatures {
		fmt.Printf("Signature:%d  %s\n", i, base58.Encode(sg))
	}

	for i, a := range trx.Message.Accounts {
		fmt.Printf("Account:%d  %s\n", i, a.ToBase58())
	}
}

func TestCreateRawTransaction(t *testing.T) {
	type args struct {
		param CreateRawTransactionParam
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			args: args{
				param: CreateRawTransactionParam{
					Instructions: []Instruction{
						{
							ProgramID: common.PublicKeyFromString("11111111111111111111111111111111"),
							Accounts: []AccountMeta{
								{
									PubKey:     common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
									IsSigner:   true,
									IsWritable: true,
								},
								{
									PubKey:     common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"),
									IsSigner:   false,
									IsWritable: true,
								},
							},
							Data: []byte{2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
						},
					},
					Signers: []Account{
						AccountFromPrivateKeyBytes([]byte{220, 190, 97, 243, 86, 180, 6, 192, 121, 120, 30, 246, 134, 81, 46, 27, 181, 181, 148, 200, 182, 184, 13, 124, 51, 186, 141, 11, 125, 116, 9, 203, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240}),
					},
					FeePayer:        common.PublicKey{},
					RecentBlockHash: "FwRYtTPRk5N4wUeP87rTw9kQVSwigB6kbikGzzeCMrW5",
				},
			},
			want:    []byte{1, 189, 98, 67, 19, 102, 99, 124, 234, 70, 209, 28, 10, 33, 66, 167, 162, 222, 122, 16, 68, 248, 129, 46, 111, 221, 255, 40, 40, 236, 84, 233, 213, 234, 185, 235, 222, 155, 204, 139, 164, 184, 155, 32, 54, 151, 73, 235, 65, 200, 76, 127, 111, 244, 72, 183, 208, 21, 247, 114, 176, 181, 21, 77, 8, 1, 0, 1, 3, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 221, 244, 189, 59, 8, 252, 7, 91, 129, 169, 22, 151, 32, 104, 208, 131, 64, 75, 232, 201, 77, 13, 187, 220, 103, 232, 190, 100, 35, 210, 17, 42, 1, 2, 2, 0, 1, 12, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			args: args{
				param: CreateRawTransactionParam{
					Instructions: []Instruction{
						{
							ProgramID: common.PublicKeyFromString("11111111111111111111111111111111"),
							Accounts: []AccountMeta{
								{
									PubKey:     common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
									IsSigner:   true,
									IsWritable: true,
								},
								{
									PubKey:     common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"),
									IsSigner:   false,
									IsWritable: true,
								},
							},
							Data: []byte{2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
						},
						{
							ProgramID: common.PublicKeyFromString("11111111111111111111111111111111"),
							Accounts: []AccountMeta{
								{
									PubKey:     common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"),
									IsSigner:   true,
									IsWritable: true,
								},
								{
									PubKey:     common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
									IsSigner:   false,
									IsWritable: true,
								},
							},
							Data: []byte{2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
						},
					},
					Signers: []Account{
						AccountFromPrivateKeyBytes([]byte{220, 190, 97, 243, 86, 180, 6, 192, 121, 120, 30, 246, 134, 81, 46, 27, 181, 181, 148, 200, 182, 184, 13, 124, 51, 186, 141, 11, 125, 116, 9, 203, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240}),
						AccountFromPrivateKeyBytes([]byte{55, 197, 194, 189, 188, 226, 127, 64, 68, 154, 221, 208, 200, 63, 127, 189, 249, 107, 106, 53, 74, 225, 149, 73, 111, 6, 153, 152, 62, 77, 118, 242, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142}),
					},
					FeePayer:        common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
					RecentBlockHash: "9qERNBLXzCqchyfquh2DjUT21xsLym6ynZPRh9TZbEiq",
				},
			},
			want:    []byte{2, 74, 231, 188, 191, 144, 39, 14, 161, 169, 155, 174, 83, 136, 177, 49, 105, 154, 137, 23, 153, 145, 47, 130, 208, 246, 195, 244, 141, 52, 228, 21, 190, 130, 99, 162, 145, 30, 133, 140, 2, 103, 40, 95, 141, 116, 111, 249, 205, 59, 137, 56, 204, 67, 132, 148, 152, 74, 69, 48, 200, 227, 0, 156, 8, 33, 150, 49, 151, 221, 70, 119, 149, 120, 244, 227, 186, 179, 109, 146, 176, 20, 58, 224, 180, 254, 64, 210, 181, 208, 226, 151, 52, 192, 198, 242, 20, 184, 23, 238, 214, 165, 140, 56, 190, 100, 122, 29, 216, 79, 196, 144, 239, 203, 64, 106, 255, 216, 27, 153, 242, 78, 154, 235, 204, 72, 58, 227, 3, 2, 0, 1, 3, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 131, 56, 140, 99, 57, 71, 67, 79, 102, 217, 86, 239, 231, 34, 85, 48, 147, 87, 18, 236, 176, 227, 54, 6, 201, 50, 117, 164, 3, 220, 147, 222, 2, 2, 2, 0, 1, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
		{
			args: args{
				param: CreateRawTransactionParam{
					Instructions: []Instruction{
						{
							ProgramID: common.PublicKeyFromString("11111111111111111111111111111111"),
							Accounts: []AccountMeta{
								{
									PubKey:     common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
									IsSigner:   true,
									IsWritable: true,
								},
								{
									PubKey:     common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"),
									IsSigner:   false,
									IsWritable: true,
								},
							},
							Data: []byte{2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
						},
						{
							ProgramID: common.PublicKeyFromString("11111111111111111111111111111111"),
							Accounts: []AccountMeta{
								{
									PubKey:     common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"),
									IsSigner:   true,
									IsWritable: true,
								},
								{
									PubKey:     common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
									IsSigner:   false,
									IsWritable: true,
								},
							},
							Data: []byte{2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
						},
					},
					Signers: []Account{
						AccountFromPrivateKeyBytes([]byte{55, 197, 194, 189, 188, 226, 127, 64, 68, 154, 221, 208, 200, 63, 127, 189, 249, 107, 106, 53, 74, 225, 149, 73, 111, 6, 153, 152, 62, 77, 118, 242, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142}),
						AccountFromPrivateKeyBytes([]byte{220, 190, 97, 243, 86, 180, 6, 192, 121, 120, 30, 246, 134, 81, 46, 27, 181, 181, 148, 200, 182, 184, 13, 124, 51, 186, 141, 11, 125, 116, 9, 203, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240}),
					},
					FeePayer:        common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
					RecentBlockHash: "9qERNBLXzCqchyfquh2DjUT21xsLym6ynZPRh9TZbEiq",
				},
			},
			want:    []byte{2, 74, 231, 188, 191, 144, 39, 14, 161, 169, 155, 174, 83, 136, 177, 49, 105, 154, 137, 23, 153, 145, 47, 130, 208, 246, 195, 244, 141, 52, 228, 21, 190, 130, 99, 162, 145, 30, 133, 140, 2, 103, 40, 95, 141, 116, 111, 249, 205, 59, 137, 56, 204, 67, 132, 148, 152, 74, 69, 48, 200, 227, 0, 156, 8, 33, 150, 49, 151, 221, 70, 119, 149, 120, 244, 227, 186, 179, 109, 146, 176, 20, 58, 224, 180, 254, 64, 210, 181, 208, 226, 151, 52, 192, 198, 242, 20, 184, 23, 238, 214, 165, 140, 56, 190, 100, 122, 29, 216, 79, 196, 144, 239, 203, 64, 106, 255, 216, 27, 153, 242, 78, 154, 235, 204, 72, 58, 227, 3, 2, 0, 1, 3, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 131, 56, 140, 99, 57, 71, 67, 79, 102, 217, 86, 239, 231, 34, 85, 48, 147, 87, 18, 236, 176, 227, 54, 6, 201, 50, 117, 164, 3, 220, 147, 222, 2, 2, 2, 0, 1, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateRawTransaction(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateRawTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRawTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransactionDeserialize(t *testing.T) {
	type args struct {
		tx []byte
	}
	tests := []struct {
		name    string
		args    args
		want    Transaction
		wantErr bool
	}{
		{
			args: args{
				tx: []byte{1, 189, 98, 67, 19, 102, 99, 124, 234, 70, 209, 28, 10, 33, 66, 167, 162, 222, 122, 16, 68, 248, 129, 46, 111, 221, 255, 40, 40, 236, 84, 233, 213, 234, 185, 235, 222, 155, 204, 139, 164, 184, 155, 32, 54, 151, 73, 235, 65, 200, 76, 127, 111, 244, 72, 183, 208, 21, 247, 114, 176, 181, 21, 77, 8, 1, 0, 1, 3, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 221, 244, 189, 59, 8, 252, 7, 91, 129, 169, 22, 151, 32, 104, 208, 131, 64, 75, 232, 201, 77, 13, 187, 220, 103, 232, 190, 100, 35, 210, 17, 42, 1, 2, 2, 0, 1, 12, 2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
			},
			want: Transaction{
				Signatures: []Signature{[]byte{189, 98, 67, 19, 102, 99, 124, 234, 70, 209, 28, 10, 33, 66, 167, 162, 222, 122, 16, 68, 248, 129, 46, 111, 221, 255, 40, 40, 236, 84, 233, 213, 234, 185, 235, 222, 155, 204, 139, 164, 184, 155, 32, 54, 151, 73, 235, 65, 200, 76, 127, 111, 244, 72, 183, 208, 21, 247, 114, 176, 181, 21, 77, 8}},
				Message: Message{
					Header: MessageHeader{
						NumRequireSignatures:        1,
						NumReadonlySignedAccounts:   0,
						NumReadonlyUnsignedAccounts: 1,
					},
					Accounts: []common.PublicKey{
						common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
						common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"),
						common.SystemProgramID,
					},
					RecentBlockHash: "FwRYtTPRk5N4wUeP87rTw9kQVSwigB6kbikGzzeCMrW5",
					Instructions: []CompiledInstruction{
						{
							ProgramIDIndex: 2,
							Accounts:       []int{0, 1},
							Data:           []byte{2, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			args: args{
				tx: []byte{2, 74, 231, 188, 191, 144, 39, 14, 161, 169, 155, 174, 83, 136, 177, 49, 105, 154, 137, 23, 153, 145, 47, 130, 208, 246, 195, 244, 141, 52, 228, 21, 190, 130, 99, 162, 145, 30, 133, 140, 2, 103, 40, 95, 141, 116, 111, 249, 205, 59, 137, 56, 204, 67, 132, 148, 152, 74, 69, 48, 200, 227, 0, 156, 8, 33, 150, 49, 151, 221, 70, 119, 149, 120, 244, 227, 186, 179, 109, 146, 176, 20, 58, 224, 180, 254, 64, 210, 181, 208, 226, 151, 52, 192, 198, 242, 20, 184, 23, 238, 214, 165, 140, 56, 190, 100, 122, 29, 216, 79, 196, 144, 239, 203, 64, 106, 255, 216, 27, 153, 242, 78, 154, 235, 204, 72, 58, 227, 3, 2, 0, 1, 3, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 131, 56, 140, 99, 57, 71, 67, 79, 102, 217, 86, 239, 231, 34, 85, 48, 147, 87, 18, 236, 176, 227, 54, 6, 201, 50, 117, 164, 3, 220, 147, 222, 2, 2, 2, 0, 1, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
			},
			want: Transaction{
				Signatures: []Signature{
					[]byte{74, 231, 188, 191, 144, 39, 14, 161, 169, 155, 174, 83, 136, 177, 49, 105, 154, 137, 23, 153, 145, 47, 130, 208, 246, 195, 244, 141, 52, 228, 21, 190, 130, 99, 162, 145, 30, 133, 140, 2, 103, 40, 95, 141, 116, 111, 249, 205, 59, 137, 56, 204, 67, 132, 148, 152, 74, 69, 48, 200, 227, 0, 156, 8},
					[]byte{33, 150, 49, 151, 221, 70, 119, 149, 120, 244, 227, 186, 179, 109, 146, 176, 20, 58, 224, 180, 254, 64, 210, 181, 208, 226, 151, 52, 192, 198, 242, 20, 184, 23, 238, 214, 165, 140, 56, 190, 100, 122, 29, 216, 79, 196, 144, 239, 203, 64, 106, 255, 216, 27, 153, 242, 78, 154, 235, 204, 72, 58, 227, 3},
				},
				Message: Message{
					Header: MessageHeader{
						NumRequireSignatures:        2,
						NumReadonlySignedAccounts:   0,
						NumReadonlyUnsignedAccounts: 1,
					},
					Accounts: []common.PublicKey{
						common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"),
						common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"),
						common.SystemProgramID,
					},
					RecentBlockHash: "9qERNBLXzCqchyfquh2DjUT21xsLym6ynZPRh9TZbEiq",
					Instructions: []CompiledInstruction{
						{
							ProgramIDIndex: 2,
							Accounts:       []int{0, 1},
							Data:           []byte{2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
						},
						{
							ProgramIDIndex: 2,
							Accounts:       []int{1, 0},
							Data:           []byte{2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransactionDeserialize(tt.args.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransactionDeserialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransactionDeserialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTransaction(t *testing.T) {
	type args struct {
		message        Message
		signaturePairs map[common.PublicKey]Signature
	}
	tests := []struct {
		name    string
		args    args
		want    Transaction
		wantErr bool
	}{
		{
			args: args{
				message: MustMessageDeserialize([]byte{2, 0, 1, 3, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 131, 56, 140, 99, 57, 71, 67, 79, 102, 217, 86, 239, 231, 34, 85, 48, 147, 87, 18, 236, 176, 227, 54, 6, 201, 50, 117, 164, 3, 220, 147, 222, 2, 2, 2, 0, 1, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}),
				signaturePairs: map[common.PublicKey]Signature{
					common.PublicKeyFromString("A4iUVr5KjmsLymUcv4eSKPedUtoaBceiPeGipKMYc69b"): {33, 150, 49, 151, 221, 70, 119, 149, 120, 244, 227, 186, 179, 109, 146, 176, 20, 58, 224, 180, 254, 64, 210, 181, 208, 226, 151, 52, 192, 198, 242, 20, 184, 23, 238, 214, 165, 140, 56, 190, 100, 122, 29, 216, 79, 196, 144, 239, 203, 64, 106, 255, 216, 27, 153, 242, 78, 154, 235, 204, 72, 58, 227, 3},
					common.PublicKeyFromString("EvN4kgKmCmYzdbd5kL8Q8YgkUW5RoqMTpBczrfLExtx7"): {74, 231, 188, 191, 144, 39, 14, 161, 169, 155, 174, 83, 136, 177, 49, 105, 154, 137, 23, 153, 145, 47, 130, 208, 246, 195, 244, 141, 52, 228, 21, 190, 130, 99, 162, 145, 30, 133, 140, 2, 103, 40, 95, 141, 116, 111, 249, 205, 59, 137, 56, 204, 67, 132, 148, 152, 74, 69, 48, 200, 227, 0, 156, 8},
				},
			},
			want:    MustTransactionDeserialize([]byte{2, 74, 231, 188, 191, 144, 39, 14, 161, 169, 155, 174, 83, 136, 177, 49, 105, 154, 137, 23, 153, 145, 47, 130, 208, 246, 195, 244, 141, 52, 228, 21, 190, 130, 99, 162, 145, 30, 133, 140, 2, 103, 40, 95, 141, 116, 111, 249, 205, 59, 137, 56, 204, 67, 132, 148, 152, 74, 69, 48, 200, 227, 0, 156, 8, 33, 150, 49, 151, 221, 70, 119, 149, 120, 244, 227, 186, 179, 109, 146, 176, 20, 58, 224, 180, 254, 64, 210, 181, 208, 226, 151, 52, 192, 198, 242, 20, 184, 23, 238, 214, 165, 140, 56, 190, 100, 122, 29, 216, 79, 196, 144, 239, 203, 64, 106, 255, 216, 27, 153, 242, 78, 154, 235, 204, 72, 58, 227, 3, 2, 0, 1, 3, 206, 211, 135, 230, 195, 111, 87, 254, 147, 239, 143, 81, 110, 159, 49, 140, 109, 137, 224, 197, 24, 49, 223, 61, 123, 8, 78, 109, 110, 136, 228, 240, 134, 172, 209, 213, 227, 137, 61, 108, 116, 171, 205, 124, 54, 68, 61, 110, 80, 31, 240, 117, 108, 137, 97, 222, 38, 242, 68, 156, 27, 65, 29, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 131, 56, 140, 99, 57, 71, 67, 79, 102, 217, 86, 239, 231, 34, 85, 48, 147, 87, 18, 236, 176, 227, 54, 6, 201, 50, 117, 164, 3, 220, 147, 222, 2, 2, 2, 0, 1, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 2, 2, 1, 0, 12, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTransaction(tt.args.message, tt.args.signaturePairs)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
