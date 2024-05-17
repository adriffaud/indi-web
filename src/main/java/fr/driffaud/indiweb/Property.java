package fr.driffaud.indiweb;

public record Property(String device, String group, PropertyType type, String name,
                       String label, String state, String perm) {
}